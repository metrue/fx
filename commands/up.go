package commands

import (
	"context"
	"io/ioutil"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/common"
	"github.com/metrue/fx/pkg/client"
	"github.com/metrue/fx/pkg/utils"
)

func InvokeUpRequest(address string, functions []string) (*api.UpResponse, error) {
	var funcList []*api.FunctionMeta
	for _, function := range functions {
		data, err := ioutil.ReadFile(function)
		if err != nil {
			return nil, err
		}

		funcMeta := &api.FunctionMeta{
			Lang:    utils.GetLangFromFileName(function),
			Path:    function,
			Content: string(data),
		}
		funcList = append(funcList, funcMeta)
	}

	req := &api.UpRequest{
		Functions: funcList,
	}

	client, conn, err := client.NewClient(address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	ctx := context.Background()
	res, err := client.Up(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func Up(address string, functions []string) error {
	res, err := InvokeUpRequest(address, functions)
	if err != nil {
		return err
	}
	common.HandleUpResult(res.Instances)

	return nil
}
