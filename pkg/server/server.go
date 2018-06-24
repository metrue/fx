package server

import (
	"context"
	"net"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/api/service"

	"google.golang.org/grpc"
)

type Fx struct {
	server *grpc.Server
	listen net.Listener
}

func NewFxServiceServer(uri string) *Fx {
	listen, err := net.Listen("tcp", uri)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	s := &Fx{
		server: server,
		listen: listen,
	}
	api.RegisterFxServiceServer(server, s)
	return s
}

func (f *Fx) Start() error {
	return f.server.Serve(f.listen)
}

func (f *Fx) Stop() {
	if f.server == nil {
		return
	}
	f.server.Stop()
	f.server = nil
}

func (f *Fx) Up(ctx context.Context, msg *api.UpRequest) (*api.UpResponse, error) {
	return service.Up(ctx, msg)
}

func (f *Fx) Down(ctx context.Context, msg *api.DownRequest) (*api.DownResponse, error) {
	return service.Down(ctx, msg)
}

func (f *Fx) List(ctx context.Context, msg *api.ListRequest) (*api.ListResponse, error) {
	return service.List(ctx, msg)
}

func (f *Fx) Ping(ctx context.Context, msg *api.PingRequest) (*api.PingResponse, error) {
	return service.Ping(ctx, msg)
}
