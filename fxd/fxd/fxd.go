package main

import (
	"context"
	"fmt"
	"log"
	"net"

	dockerSDK "github.com/metrue/fx/container_runtimes/docker/sdk"
	"github.com/metrue/fx/driver"
	dockerDriver "github.com/metrue/fx/driver/docker"

	"github.com/metrue/fx/api"
	"google.golang.org/grpc"
)

type server struct {
	driver driver.Driver
}

func (s *server) ListService(ctx context.Context, in *api.ListServiceRequest) (*api.ListServiceResponse, error) {
	filter := in.GetFilter()
	fmt.Println("++++++++")
	fmt.Println(filter)
	fmt.Println("++++++++")
	services, err := s.driver.List(ctx, filter.GetName())
	if err != nil {
		return nil, err
	}
	fmt.Println(services)
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

	docker, err := dockerSDK.CreateClient(context.Background())
	if err != nil {
		log.Fatalf("failed to create a docker cli: %v", err)
	}

	driver := dockerDriver.New(dockerDriver.Options{
		DockerClient: docker,
	})
	api.RegisterAPIServer(s, &server{driver: driver})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
