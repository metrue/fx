package commands

import (
	"context"
	"io/ioutil"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/common"
	"github.com/metrue/fx/pkg/client"
	"github.com/metrue/fx/pkg/utils"
	"github.com/pkg/errors"
)

var (
	ReadSourceError = errors.New("Could not read function source")
	UpFunctionError = errors.New("Could not up function")
)

// Up starts the functions specified in flags
func Up(address string, functions []string) error {
	var funcList []*api.FunctionMeta
	for _, function := range functions {
		data, err := ioutil.ReadFile(function)
		if err != nil {
			return ReadSourceError
		}

		funcMeta := &api.FunctionMeta{
			Lang:    utils.GetLangFromFileName(function),
			Path:    function,
			Content: string(data),
		}
		funcList = append(funcList, funcMeta)
	}

	client, conn, err := client.NewClient(address)
	if err != nil {
		return NewClientError
	}

	defer conn.Close()

	ctx := context.Background()
	req := &api.UpRequest{
		Functions: funcList,
	}
	res, err := client.Up(ctx, req)
	if err != nil {
		return UpFunctionError
	}

	common.HandleUpResult(res.Instances)
	return nil
}
