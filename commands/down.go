package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/common"
)

// Down invoke the removal of one or more functions
func Down() {
	option := "down"
	nArgs := len(os.Args)
	args, flagSet := common.SetupFlags(option)
	if nArgs == 2 {
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
	req := &api.DownRequest{
		ID: functions,
	}
	res, err := client.Down(ctx, req)

	if err != nil {
		panic(err)
	}

	fmt.Println(DownMessage(res.Instances))
}
