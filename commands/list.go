package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/common"
)

// List lists all running function services
func List() {
	option := "list"
	nArgs := len(os.Args)
	args, flagSet := common.SetupFlags(option)
	if nArgs < 2 {
		common.FlagsAndExit(flagSet)
	}
	functions, address := common.ParseArgs(
		option,
		os.Args[2:],
		args,
		flagSet,
	)

	client, conn, err := api.NewClient(address)
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	ctx := context.Background()
	req := &api.ListRequest{
		ID: functions,
	}
	res, err := client.List(ctx, req)

	if err != nil {
		panic(err)
	}

	fmt.Println(ListMessage(res.Instances))
}
