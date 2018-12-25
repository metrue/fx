package commands

import (
	"context"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/pkg/client"
	"github.com/metrue/fx/pkg/utils"
)

func Status(address string) error {
	client, conn, err := client.NewClient(address)
	if err != nil {
		return err
	}
	defer conn.Close()

	ctx := context.Background()
	req := &api.PingRequest{}
	_, err = client.Ping(ctx, req)
	if err != nil {
		return err
	}
	return utils.OutputJSON(map[string]string{
		"status": "ok",
		"server": address,
	})
}
