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
	upReq := &api.UpRequest{
		Functions: []*api.FunctionMeta{
			&api.FunctionMeta{
				Lang:    "node",
				Content: "module.exports = () => { return \"foo\"; }",
				Path:    "./",
			},
		},
	}
	_, err = client.Up(ctx, upReq)
	if err != nil {
		t.Fatal(err)
		return
	}

	ctx = context.Background()
	listReq := &api.ListRequest{}
	listRes, err := client.List(ctx, listReq)
	if err != nil {
		t.Fatal(err)
		return
	}

	for i := 0; i < listRes.Instances; i++ {
		instance := listRes.Instances[i]
	}

	conn.Close()
	service.Stop()
}
