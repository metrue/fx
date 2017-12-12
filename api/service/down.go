package service

import (
	"github.com/metrue/fx/api"
	"github.com/metrue/fx/handlers"
)

//downTask wrap a DownMsgMeta and an error from its processing
type downTask struct {
	Val *api.DownMsgMeta
	Err error
}

//newDownTask initialize a new downTask
func newDownTask(val *api.DownMsgMeta, err error) downTask {
	return downTask{
		Val: val,
		Err: err,
	}
}

//Down handle function removal requests
func Down(req *api.DownRequest) (*api.DownResponse, error) {

	containers := handlers.List(req.ID...)
	count := len(containers)
	results := make(chan downTask, count)

	for _, container := range containers {
		go func() {
			results <- newDownTask(handlers.Down(container.ID[:10], container.Image))
		}()
	}

	// collect down result
	var downs []*api.DownMsgMeta
	for result := range results {
		downResult := result.Val
		if result.Err != nil {
			downResult.Error = result.Err.Error()
		}
		downs = append(downs, downResult)
		if len(downs) == count {
			close(results)
		}
	}

	return &api.DownResponse{
		Instances: downs,
	}, nil
}
