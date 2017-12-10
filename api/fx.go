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

//Start the gRPC server
func Start(uri string) error {

	if server != nil {
		return nil
	}

	listen, err := net.Listen("tcp", uri)
	if err != nil {
		return err
	}
	server = grpc.NewServer()
	RegisterFxServiceServer(server, newFxService())
	server.Serve(listen)
	return nil
}

//Stop the server
func Stop() {
	if server == nil {
		return
	}
	server.Stop()
}

func (f *fx) Up(ctx context.Context, msg *UpRequest) (*UpResponse, error) {
	return nil, errors.New("Not implemented")
}

func (f *fx) Down(ctx context.Context, msg *DownRequest) (*DownResponse, error) {
	return nil, errors.New("Not implemented")
}

func (f *fx) List(ctx context.Context, msg *ListRequest) (*ListResponse, error) {
	return nil, errors.New("Not implemented")
}
