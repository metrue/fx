package commands

import (
	"context"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/pkg/client"
	"github.com/metrue/fx/pkg/utils"
)

// List lists all running function services
func List(address string, functions []string) error {
	client, conn, err := client.NewClient(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx := context.Background()
	req := &api.ListRequest{
		ID: functions,
	}

	res, err := client.List(ctx, req)
	if err != nil {
		return err
	}

	return utils.OutputJSON(res)
}
