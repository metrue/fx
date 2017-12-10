package api

import (
	"github.com/metrue/fx/common"
	"github.com/metrue/fx/handlers"
	Message "github.com/metrue/fx/message"
)

func Up(req *UpRequest) (*UpResponse, error) {

	funcList := req.Functions

	count := len(funcList)
	upResultCh := make(chan Message.UpMsgMeta, count)
	for _, funcMeta := range funcList {
		//TODO use only one type avoiding conversion where possible
		meta := common.FunctionMeta{
			Content: funcMeta.Content,
			Lang:    funcMeta.Lang,
			Path:    funcMeta.Path,
		}
		go handlers.Up(meta, upResultCh)
	}

	// collect down result
	var ups []*UpMsgMeta
	for upResult := range upResultCh {
		//TODO use only one type avoiding conversion where possible
		ups = append(ups, &UpMsgMeta{
			FunctionSource: upResult.FunctionSource,
			LocalAddress:   upResult.LocalAddress,
			RemoteAddress:  upResult.RemoteAddress,
		})
		if len(ups) == count {
			close(upResultCh)
		}
	}

	res := &UpResponse{
		Instances: ups,
	}

	return res, nil
}
