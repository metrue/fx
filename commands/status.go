package commands

import (
	"context"
	"fmt"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/pkg/client"
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

	info := fmt.Sprintf(`
---
status: ok ):
fx server: %s
---
	`, address)
	fmt.Println(info)

	return nil
}
