package service

import (
	"github.com/metrue/fx/api"
	"github.com/metrue/fx/handlers"
)

//upTask wrap the UpMsgMeta and an error from processing
type upTask struct {
	Val *api.UpMsgMeta
	Err error
}

//newUpTask initialize a new upTask
func newUpTask(val *api.UpMsgMeta, err error) upTask {
	return upTask{
		Val: val,
		Err: err,
	}
}

//Up deploy and run functions
func UpDockerfile(req *api.UpDockerfileRequest) (*api.UpResponse, error) {
	dockerfile = req.dockerfile;

	funcList := req.Functions;

	count := len(funcList)
	results := make(chan upTask, count)

	for _, meta := range funcList {
		go func(funcMeta *api.FunctionMeta) {
			results <- newUpTask(handlers.Up(*funcMeta))
		}(meta)
	}

	// collect up results
	var ups []*api.UpMsgMeta
	for result := range results {
		upResult := result.Val
		if result.Err != nil {
			upResult = &api.UpMsgMeta{
				Error: result.Err.Error(),
			}
		}
		ups = append(ups, upResult)
		if len(ups) == count {
			close(results)
		}
	}

	return &api.UpResponse{
		Instances: ups,
	}, nil
}
