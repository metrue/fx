package service_test

import (
	"context"
	"testing"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/api/service"
)

func TestDown(t *testing.T) {

	runServer(t)

	client, conn, err := api.NewClient(grpcEndpoint)
	if err != nil {
		t.Fatal(err)
		return
	}

	ctx := context.Background()
	req := &api.UpRequest{}
	res, err = client.Up(ctx, req)
	if err != nil {
		t.Fatal(err)
		return
	}

	ctx := context.Background()
	req := &api.DownRequest{}
	_, err = client.List(ctx, req)
	if err != nil {
		t.Fatal(err)
		return
	}

	conn.Close()
	service.Stop()
}
