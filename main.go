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
	"github.com/hkwi/h2c"
	"google.golang.org/grpc"
)

type doerServer struct {
}

func (s *doerServer) DoIt(ctx context.Context, c *pb.Command) (*pb.Response, error) {
	resp := fmt.Sprintf("Did: %s", c.Thing)
	log.Printf("RPC: %s\n DONE!", c.Thing)
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

		log.Printf("STREAM: %s", in)

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
	log.Printf("HTTP/%d %s: %s\n", r.ProtoMajor, r.Method, in)
	fmt.Fprintf(w, "You %s for %q\n\n", r.Method, r.URL)
	fmt.Fprintf(w, "Got: %s\n", in)
}

func newServer() *doerServer {
	s := &doerServer{}
	return s
}

type grpcAdapter struct {
	grpcServer http.Handler
}

func (g grpcAdapter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("Got %s on %d with %s", r.Method, r.ProtoMajor, r.Header.Get("Content-Type"))
	if r.ProtoMajor == 2 && strings.HasPrefix(
		r.Header.Get("Content-Type"), "application/grpc") {
		g.grpcServer.ServeHTTP(w, r)
	} else {
		http.DefaultServeMux.ServeHTTP(w, r)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	fmt.Printf("Listening on %s\n", addr)
	grpcServer := grpc.NewServer()
	doer := newServer()
	pb.RegisterDoerServer(grpcServer, doer)
	http.Handle("/", doer)

	h2g := grpcAdapter{grpcServer: grpcServer}
	noTls := h2c.Server{Handler: h2g}
	log.Fatal(http.ListenAndServe(addr, noTls))
}
