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
		go handlers.Up(funcMeta, upResultCh)
	}

	// collect down result
	var ups []*api.UpMsgMeta
	for upResult := range upResultCh {
		ups = append(ups, upResult)
		if len(ups) == count {
			close(upResultCh)
		}
	}

	res := &api.UpResponse{
		Instances: ups,
	}

	return res, nil
}
