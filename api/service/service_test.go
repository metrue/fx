package service

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/metrue/fx/api"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

const grpcEndpoint = "localhost:5001"

var client api.FxServiceClient

func startServer() {
	go func() {
		err := Start(grpcEndpoint)
		if err != nil {
			panic(err)
		}
	}()
	//wait for the service to start
	time.Sleep((time.Millisecond * 2000))
}

func stopServer(conn *grpc.ClientConn) {
	conn.Close()
	Stop()
	//wait for the service to start
	time.Sleep((time.Millisecond * 2000))
}

func TestPingService(t *testing.T) {
	ctx := context.Background()
	req := &api.PingRequest{}
	res, err := client.Ping(ctx, req)
	assert.Nil(t, err)
	assert.Equal(t, res, &api.PingResponse{Status: "pong"})
}

func TestListService(t *testing.T) {
	ctx := context.Background()
	req := &api.ListRequest{}
	_, err := client.List(ctx, req)
	assert.Nil(t, err)
}

func TestUpService(t *testing.T) {
	assert.Nil(t, nil)
}

func TestDownService(t *testing.T) {
	assert.Nil(t, nil)
}

func TestMain(m *testing.M) {
	startServer()

	cli, conn, err := api.NewClient(grpcEndpoint)
	if err != nil {
		panic(err)
	}
	client = cli
	defer stopServer(conn)

	os.Exit(m.Run())

}
