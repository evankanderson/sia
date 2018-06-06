# build stage
FROM golang:1.10 AS build-env
RUN mkdir -p /go/src/github.com/evankanderson/sia
WORKDIR /go/src/github.com/evankanderson/sia
ADD . .
RUN go get -v ./...
# RUN go install -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w' -o doerapp

# final stage
FROM alpine
WORKDIR /app
COPY --from=build-env /go/src/github.com/evankanderson/sia/doerapp /app/doer
ENTRYPOINT ["/app/doer"]
ENV PORT=8080
# Actually $PORT
EXPOSE 8080
