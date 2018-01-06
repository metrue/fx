package commands

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/common"
	"github.com/metrue/fx/utils"
	"github.com/pkg/errors"
)

// Up starts the functions specified in flags
func Up() {
	option := "up"
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

	fmt.Println("Deploy starting...")

	channel := make(chan bool)
	defer close(channel)

	var funcList []*api.FunctionMeta
	for _, function := range functions {
		data, err := ioutil.ReadFile(function)
		if err != nil {
			err = errors.Wrap(err, "Read function content falied")
			common.HandleError(err)
		}

		funcMeta := &api.FunctionMeta{
			Lang:    utils.GetLangFromFileName(function),
			Path:    function,
			Content: string(data),
		}
		funcList = append(funcList, funcMeta)
	}

	client, conn, err := api.NewClient(address)
	if err != nil {
		err = errors.Wrap(err, "New gRPC Client failed")
		common.HandleError(err)
	}

	defer conn.Close()

	ctx := context.Background()
	req := &api.UpRequest{
		Functions: funcList,
	}
	res, err := client.Up(ctx, req)
	if err != nil {
		err = errors.Wrap(err, "up function failed")
		common.HandleError(err)
	}

	common.HandleUpResult(res.Instances)
}
