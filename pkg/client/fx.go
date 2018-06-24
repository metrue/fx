package client

import (
	"github.com/metrue/fx/api"
	"google.golang.org/grpc"
)

//NewClient return a new gRPC client with default settings
func NewClient(grpcEndpoint string) (api.FxServiceClient, *grpc.ClientConn, error) {
	var opts []grpc.DialOption

	//TODO review options
	opts = append(opts, grpc.WithInsecure())

	conn, err := grpc.Dial(grpcEndpoint, opts...)
	if err != nil {
		return nil, nil, err
	}

	client := api.NewFxServiceClient(conn)
	return client, conn, nil
}
