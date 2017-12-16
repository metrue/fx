package service

import (
	"context"
	"testing"
	"time"

	"github.com/metrue/fx/api"
	"google.golang.org/grpc"
)

const grpcEndpoint = "localhost:5001"

func runServer(t *testing.T) {
	go func() {
		err := Start(grpcEndpoint)
		if err != nil {
			t.Fatal(err)
		}
	}()
	//wait for the service to start
	time.Sleep((time.Millisecond * 500))
}

func stopServer(conn *grpc.ClientConn) {
	conn.Close()
	Stop()
	//wait for the service to start
	time.Sleep((time.Millisecond * 500))
}

func TestServer(t *testing.T) {

	runServer(t)

	client, conn, err := api.NewClient(grpcEndpoint)
	if err != nil {
		t.Fatal(err)
	}
	defer stopServer(conn)

	ctx := context.Background()
	req := &api.ListRequest{}
	_, err = client.List(ctx, req)
	if err != nil {
		t.Fatal(err)
	}

}
