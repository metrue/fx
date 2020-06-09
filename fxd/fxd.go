package main

import (
	"context"
	"log"
	"net"

	fxd "github.com/metrue/fx/fxd/proto"
	"google.golang.org/grpc"
)

type server struct {
}

func (s *server) SayHelloAgain(ctx context.Context, in *fxd.ListServiceRequest) (*fxd.ListServiceResponse, error) {
	return &fxd.ListServiceResponse{Message: "Hello again " + in.GetName()}, nil
}

const (
	port = ":8866"
)

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	fxd.RegisterFxdServer(s, &server{})
	if err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
