package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/metrue/fx/api"
)

func removeFunction(client api.FxServiceClient, functionIDs ...string) (bool, error) {

	ctx := context.Background()
	downReq := &api.DownRequest{ID: functionIDs}
	downRes, err := client.Down(ctx, downReq)

	if err != nil {
		return false, err
	}

	for _, f := range downRes.Instances {
		if f.Error != "" {
			return false, fmt.Errorf("[%s] %s", f.ContainerId, f.Error)
		}
	}

	return true, nil
}

func createFunction(client api.FxServiceClient) (*api.UpMsgMeta, error) {

	ctx := context.Background()
	upReq := &api.UpRequest{
		Functions: []*api.FunctionMeta{
			&api.FunctionMeta{
				Lang:    "node",
				Content: fmt.Sprintf("module.exports = () => { return \"foo_%d\"; }", time.Now().Unix()),
				Path:    "./",
			},
		},
	}

	upRes, err := client.Up(ctx, upReq)
	if err != nil {
		return nil, err
	}

	if len(upRes.Instances) != 1 {
		return nil, fmt.Errorf("Up response should have one instance, found %d", len(upRes.Instances))
	}

	if upRes.Instances[0].Error != "" {
		return nil, fmt.Errorf("Up error: %s", upRes.Instances[0].Error)
	}

	return upRes.Instances[0], nil
}

func TestList(t *testing.T) {

	runServer(t)

	client, conn, err := api.NewClient(grpcEndpoint)
	defer stopServer(conn)

	if err != nil {
		t.Fatal(err)
	}

	upMeta, err := createFunction(client)
	if err != nil {
		t.Fatalf("Up failed: %s\n", err.Error())
	}

	ctx := context.Background()
	listReq := &api.ListRequest{}
	listRes, err := client.List(ctx, listReq)
	if err != nil {
		t.Fatal(err)
	}

	var found bool
	var ids []string
	for _, fx := range listRes.Instances {
		if fx.FunctionID == upMeta.FunctionID {
			found = true
			ids = append(ids, fx.FunctionID)
		}
	}

	if !found {
		fmt.Printf("FAIL: function not found %s\n", upMeta.FunctionID)
		t.Fail()
	}

	if _, err := removeFunction(client, ids...); err != nil {
		t.Fatalf("Failed to remove functions: %s\n", err.Error())
	}

}
