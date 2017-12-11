package service_test

import (
	"context"
	"testing"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/api/service"
)

const grpcEndpoint = ":5001"

func runServer(t *testing.T) {
	go func() {
		err := service.Start(grpcEndpoint)
		if err != nil {
			t.Fatal(err)
			return
		}
	}()
}

func TestServer(t *testing.T) {

	runServer(t)

	client, conn, err := api.NewClient(grpcEndpoint)
	if err != nil {
		t.Fatal(err)
		return
	}

	ctx := context.Background()
	req := &api.ListRequest{}
	_, err = client.List(ctx, req)
	if err != nil {
		t.Fatal(err)
		return
	}

	conn.Close()
	service.Stop()
}
