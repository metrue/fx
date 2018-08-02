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

type CallOutput struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func InvokeCallRequest(address string, function string, params string) (*api.CallResponse, error) {
	data, err := ioutil.ReadFile(function)
	if err != nil {
		return nil, err
	}

	req := &api.CallRequest{
		Lang:    utils.GetLangFromFileName(function),
		Content: string(data),
		Params:  params,
	}

	client, conn, err := client.NewClient(address)
	if err != nil {
		fmt.Println(client, conn, err)
		return nil, err
	}
	defer conn.Close()

	ctx := context.Background()
	res, err := client.Call(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func Call(address string, function string, params string) error {
	res, err := InvokeCallRequest(address, function, params)
	if err != nil {
		common.HandleCallResult(CallOutput{Error: err.Error()})
		return nil
	}

	common.HandleCallResult(res)
	return nil
}
