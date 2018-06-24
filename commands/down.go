package commands

import (
	"context"
	"errors"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/common"
	"github.com/metrue/fx/pkg/client"
)

var (
	NewClientError    = errors.New("Could create a client")
	DownFunctionError = errors.New("Could not down function")
)

// Down invoke the removal of one or more functions
func Down(address string, functions []string) error {
	client, conn, err := client.NewClient(address)
	if err != nil {
		return NewClientError
	}
	defer conn.Close()

	ctx := context.Background()
	req := &api.DownRequest{
		ID: functions,
	}
	res, err := client.Down(ctx, req)
	if err != nil {
		return DownFunctionError
	}

	common.HandleDownResult(res.Instances)
	return nil
}
