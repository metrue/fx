package commands

import (
	"context"
	"os"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/common"
	"github.com/metrue/fx/config"
	"github.com/pkg/errors"
)

// List lists all running function services
func List() {
	option := "list"
	nArgs := len(os.Args)
	args, flagSet := common.SetupFlags(option)
	if nArgs < 2 {
		common.FlagsAndExit(flagSet)
	}
	functions, _ := common.ParseArgs(
		option,
		os.Args[2:],
		args,
		flagSet,
	)

	client, conn, err := api.NewClient(config.GrpcEndpoint)
	if err != nil {
		err = errors.Wrap(err, "New client failed")
		common.HandleError(err)
	}

	defer conn.Close()

	ctx := context.Background()
	req := &api.ListRequest{
		ID: functions,
	}
	res, err := client.List(ctx, req)
	if err != nil {
		err = errors.Wrap(err, "List deployed functions failed")
		common.HandleError(err)
	}

	common.HandleListResult(res.Instances)
}
