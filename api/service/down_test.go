package service

import (
	"context"
	"errors"
	"testing"

	"github.com/metrue/fx/api"
)

func containsFunction(client api.FxServiceClient, functionIDs ...string) (bool, int, error) {
	ctx := context.Background()
	listReq := &api.ListRequest{}
	listRes, err := client.List(ctx, listReq)
	if err != nil {
		return false, 0, err
	}
	var found int
	for _, fx := range listRes.Instances {
		for _, fxID := range functionIDs {
			if fx.FunctionID == fxID {
				found++
			}
		}
	}
	if found != len(functionIDs) {
		return false, found, errors.New("Missing functions")
	}
	return true, found, nil
}

func TestDownAll(t *testing.T) {

	runServer(t)

	client, conn, err := api.NewClient(grpcEndpoint)
	defer stopServer(conn)

	if err != nil {
		t.Fatal(err)
	}

	fx1, err := createFunction(client)
	if err != nil {
		t.Fatalf("Up failed: %s\n", err.Error())
	}
	fx2, err := createFunction(client)
	if err != nil {
		t.Fatalf("Up failed: %s\n", err.Error())
	}

	if _, _, err := containsFunction(client, fx1.FunctionID, fx2.FunctionID); err != nil {
		t.Fatalf(err.Error())
	}

	if _, err := removeFunction(client, "*"); err != nil {
		t.Fatalf("Failed to remove functions: %s\n", err.Error())
	}

	if _, _, err := containsFunction(client); err != nil {
		t.Fatalf(err.Error())
	}

}
