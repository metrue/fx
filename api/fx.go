//go:generate sh gen.sh

package api

import (
	"context"
	"errors"
	"net"

	"google.golang.org/grpc"
)

var server *grpc.Server

type fx struct{}

func newFxService() FxServiceServer {
	return new(fx)
}

//Start the gRPC server
func Start(uri string) error {

	if server != nil {
		return errors.New("gRPC uri not provided")
	}

	listen, err := net.Listen("tcp", uri)
	if err != nil {
		return err
	}
	server = grpc.NewServer()
	RegisterFxServiceServer(server, newFxService())

	return server.Serve(listen)
}

//Stop the gRPC server
func Stop() {
	if server == nil {
		return
	}
	server.Stop()
}

func (f *fx) Up(ctx context.Context, msg *UpRequest) (*UpResponse, error) {
	return Up(msg)
}

func (f *fx) Down(ctx context.Context, msg *DownRequest) (*DownResponse, error) {
	return Down(msg)
}

func (f *fx) List(ctx context.Context, msg *ListRequest) (*ListResponse, error) {
	return List(msg)
}
