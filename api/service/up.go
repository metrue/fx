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
func Up(req *api.UpRequest) (*api.UpResponse, error) {

	funcList := req.Functions

	count := len(funcList)
	results := make(chan upTask, count)

	for _, funcMeta := range funcList {
		go func() {
			results <- newUpTask(handlers.Up(*funcMeta))
		}()
	}

	// collect up results
	var ups []*api.UpMsgMeta
	for result := range results {
		upResult := result.Val
		if result.Err != nil {
			upResult.Error = result.Err.Error()
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
