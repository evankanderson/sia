package main

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	pb "github.com/evankanderson/sia/doer"
	"google.golang.org/grpc"
)

type doerServer struct {
}

func (s *doerServer) DoIt(ctx context.Context, c *pb.Command) (*pb.Response, error) {
	resp := fmt.Sprintf("Did: %s", c.Thing)
	return &pb.Response{Words: resp}, nil
}

func (s *doerServer) KeepDoing(stream pb.Doer_KeepDoingServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		words := fmt.Sprintf("Did: %s", in.Thing)
		resp := &pb.Response{Words: words}
		if err = stream.Send(resp); err != nil {
			return err
		}
	}
}

func (s *doerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	in, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Failed to read: %v\n", err)
		return
	}
	w.WriteHeader(200)
	fmt.Fprintf(w, "Did: %s\n", in)
}

func newServer() *doerServer {
	s := &doerServer{}
	return s
}

func main() {
	addr := fmt.Sprintf("localhost:%s", os.Getenv("PORT"))
	fmt.Printf("Listening on %s", addr)
	grpcServer := grpc.NewServer()
	doer := newServer()
	pb.RegisterDoerServer(grpcServer, doer)

	h2g := http.NewServeMux()
	h2g.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.HasPrefix(
			r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			http.DefaultServeMux.ServeHTTP(w, r)
		}
	})

	http.Handle("/", doer)

	log.Fatal(http.ListenAndServe(addr, h2g))
}
