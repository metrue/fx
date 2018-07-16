package commands

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/common"
	"github.com/metrue/fx/pkg/client"
	"github.com/metrue/fx/pkg/utils"
)

func Up(address string, functions []string) (string, string, error) {
	var funcList []*api.FunctionMeta
	for _, function := range functions {
		data, err := ioutil.ReadFile(function)
		if err != nil {
			return "", "", err
		}

		funcMeta := &api.FunctionMeta{
			Lang:    utils.GetLangFromFileName(function),
			Path:    function,
			Content: string(data),
		}
		fmt.Println(funcMeta)
		funcList = append(funcList, funcMeta)
	}

	client, conn, err := client.NewClient(address)
	if err != nil {
		return "", "", err
	}

	defer conn.Close()

	ctx := context.Background()
	req := &api.UpRequest{
		Functions: funcList,
	}
	res, err := client.Up(ctx, req)
	if err != nil {
		return "", "", err
	}

	common.HandleUpResult(res.Instances)
	ret := res.Instances[0]
	// TODO should reture a request-able Adress
	return ret.FunctionID, ret.LocalAddress, nil
}
