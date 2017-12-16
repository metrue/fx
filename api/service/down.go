package service

import (
	"github.com/docker/docker/api/types"
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

	// handle fx down *
	var ids []string
	if req.ID != nil && len(req.ID) > 0 {
		if req.ID[0] != "*" {
			ids = req.ID
		}
	}

	containers, err := handlers.List(ids...)
	if err != nil {
		return nil, err
	}
	count := len(containers)
	results := make(chan downTask, count)
	downResponse := &api.DownResponse{}

	if count == 0 {
		return downResponse, nil
	}

	for _, c := range containers {
		go func(container types.Container) {
			results <- newDownTask(handlers.Down(container.ID[:10], container.ImageID))
		}(c)
	}

	// collect down result
	var downs []*api.DownMsgMeta
	for result := range results {
		downResult := result.Val
		if result.Err != nil {
			downResult = &api.DownMsgMeta{
				Error: result.Err.Error(),
			}
		}
		downs = append(downs, downResult)
		if len(downs) == count {
			close(results)
		}
	}

	downResponse.Instances = downs
	return downResponse, nil
}
