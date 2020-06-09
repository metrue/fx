package main

import (
	"context"
	"log"
	"net"

	"github.com/metrue/fx/api"
	"google.golang.org/grpc"
)

type server struct {
}

func (s *server) ListService(ctx context.Context, in *api.ListServiceRequest) (*api.ListServiceResponse, error) {
	return &api.ListServiceResponse{}, nil
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
	api.RegisterAPIServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
