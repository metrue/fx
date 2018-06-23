package commands

import (
	"context"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/common"
	"github.com/pkg/errors"
)

// List lists all running function services
func List(address string, functions []string) error {
	client, conn, err := api.NewClient(address)
	if err != nil {
		return NewClientError
	}
	defer conn.Close()

	ctx := context.Background()
	req := &api.ListRequest{
		ID: functions,
	}

	res, err := client.List(ctx, req)
	if err != nil {
		return errors.New("Could not list functions")
	}

	common.HandleListResult(res.Instances)
	return nil
}
