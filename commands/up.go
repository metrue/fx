package commands

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/common"
	"github.com/metrue/fx/decider"
	"github.com/metrue/fx/utils"
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

	d := decider.NewDecider()
	upTargetType := d.GetUpTargetType(functions)
	if upTargetType == decider.UP_TARGET_IS_FUNCTION {

	} else if upTargetType == decider.UP_TARGET_IS_DOCKER_FILE {
		fmt.Println("to up dockerfile")
	}

	os.Exit()

	fmt.Println("Deploy starting...")

	channel := make(chan bool)
	defer close(channel)

	var funcList []*api.FunctionMeta
	for _, function := range functions {
		data, err := ioutil.ReadFile(function)
		if err != nil {
			panic(err)
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
		panic(err)
	}

	defer conn.Close()

	ctx := context.Background()
	req := &api.UpRequest{
		Functions: funcList,
	}
	res, err := client.Up(ctx, req)

	if err != nil {
		panic(err)
	}

	fmt.Println(UpMessage(res.Instances))
}
