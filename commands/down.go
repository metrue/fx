package commands

import (
	"context"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/common"
)

// Down invoke the removal of one or more functions
func Down(address string, functions []string) error {
	client, conn, err := api.NewClient(address)
	if err != nil {
		common.HandleError(err)
		return err
	}

	defer conn.Close()

	ctx := context.Background()
	req := &api.DownRequest{
		ID: functions,
	}
	res, err := client.Down(ctx, req)

	if err != nil {
		common.HandleError(err)
		return err
	}

	common.HandleDownResult(res.Instances)
	return nil
}
