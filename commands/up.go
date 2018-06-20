package commands

import (
	"context"
	"io/ioutil"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/common"
	"github.com/metrue/fx/utils"
	"github.com/pkg/errors"
)

// Up starts the functions specified in flags
func Up(address string, functions []string) error {
	var funcList []*api.FunctionMeta
	for _, function := range functions {
		data, err := ioutil.ReadFile(function)
		if err != nil {
			common.HandleError(err)
			return errors.Wrap(err, "Read function content falied")
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
		common.HandleError(err)
		return errors.Wrap(err, "New gRPC Client failed")
	}

	defer conn.Close()

	ctx := context.Background()
	req := &api.UpRequest{
		Functions: funcList,
	}
	res, err := client.Up(ctx, req)
	if err != nil {
		common.HandleError(err)
		return errors.Wrap(err, "up function failed")
	}

	common.HandleUpResult(res.Instances)

	return nil
}
