package commands

import (
	"context"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/common"
	"github.com/metrue/fx/pkg/client"
)

func Down(address string, functions []string) error {
	client, conn, err := client.NewClient(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx := context.Background()
	req := &api.DownRequest{
		ID: functions,
	}
	res, err := client.Down(ctx, req)
	if err != nil {
		return err
	}

	common.HandleDownResult(res.Instances)
	return nil
}
