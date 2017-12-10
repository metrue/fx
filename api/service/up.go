package service

import (
	"github.com/metrue/fx/api"
	"github.com/metrue/fx/handlers"
)

func Up(req *api.UpRequest) (*api.UpResponse, error) {

	funcList := req.Functions

	count := len(funcList)
	upResultCh := make(chan api.UpMsgMeta, count)
	for _, funcMeta := range funcList {
		//TODO use only one type avoiding conversion where possible
		meta := api.FunctionMeta{
			Content: funcMeta.Content,
			Lang:    funcMeta.Lang,
			Path:    funcMeta.Path,
		}
		go handlers.Up(meta, upResultCh)
	}

	// collect down result
	var ups []*api.UpMsgMeta
	for upResult := range upResultCh {
		//TODO use only one type avoiding conversion where possible
		ups = append(ups, &api.UpMsgMeta{
			FunctionSource: upResult.FunctionSource,
			LocalAddress:   upResult.LocalAddress,
			RemoteAddress:  upResult.RemoteAddress,
		})
		if len(ups) == count {
			close(upResultCh)
		}
	}

	res := &api.UpResponse{
		Instances: ups,
	}

	return res, nil
}
