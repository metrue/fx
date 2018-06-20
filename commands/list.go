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
		common.HandleError(err)
		return errors.Wrap(err, "New client failed")
	}
	defer conn.Close()

	ctx := context.Background()
	req := &api.ListRequest{
		ID: functions,
	}

	res, err := client.List(ctx, req)
	if err != nil {
		common.HandleError(err)
		return errors.Wrap(err, "List deployed functions failed")
	}

	common.HandleListResult(res.Instances)
	return nil
}
