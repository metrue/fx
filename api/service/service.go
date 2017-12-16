package service

import (
	"errors"
	"net"

	"golang.org/x/net/context"

	"github.com/metrue/fx/api"
	"google.golang.org/grpc"
)

var server *grpc.Server

type fx struct{}

func newFxService() api.FxServiceServer {
	return new(fx)
}

//Start the gRPC server
func Start(uri string) error {

	if uri == "" {
		return errors.New("gRPC uri not provided")
	}

	listen, err := net.Listen("tcp", uri)
	if err != nil {
		return err
	}
	server = grpc.NewServer()
	api.RegisterFxServiceServer(server, newFxService())

	return server.Serve(listen)
}

//Stop the gRPC server
func Stop() {
	if server == nil {
		return
	}
	server.Stop()
	server = nil
}

func (f *fx) Up(ctx context.Context, msg *api.UpRequest) (*api.UpResponse, error) {
	return Up(msg)
}

func (f *fx) Down(ctx context.Context, msg *api.DownRequest) (*api.DownResponse, error) {
	return Down(msg)
}

func (f *fx) List(ctx context.Context, msg *api.ListRequest) (*api.ListResponse, error) {
	return List(msg)
}
