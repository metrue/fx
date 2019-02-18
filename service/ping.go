package service

import (
	"github.com/metrue/fx/api"
	"golang.org/x/net/context"
)

func Ping(ctx context.Context, msg *api.PingRequest) (*api.PingResponse, error) {
	return &api.PingResponse{Status: "pong"}, nil
}
