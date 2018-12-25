package commands

import (
	"context"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/pkg/client"
	"github.com/metrue/fx/pkg/utils"
)

func InvokeDownRequest(address string, functions []string) (*api.DownResponse, error) {
	client, conn, err := client.NewClient(address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ctx := context.Background()
	req := &api.DownRequest{
		ID: functions,
	}
	res, err := client.Down(ctx, req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func Down(address string, functions []string) error {
	res, err := InvokeDownRequest(address, functions)
	if err != nil {
		return err
	}
	return utils.OutputJSON(res)
}
