package service

import (
	"github.com/metrue/fx/api"
	"github.com/metrue/fx/handlers"
)

func Down(req *api.DownRequest) (*api.DownResponse, error) {

	containers := handlers.List(req.ID...)
	count := len(containers)
	downResultCh := make(chan api.DownMsgMeta, count)

	for _, container := range containers {
		go handlers.Down(container.ID[:10], container.Image, downResultCh)
	}

	// collect down result
	var downs []*api.DownMsgMeta
	for downResult := range downResultCh {
		//TODO use only one type avoiding conversion where possible
		downs = append(downs, &api.DownMsgMeta{
			ContainerId:     downResult.ContainerId,
			ContainerStatus: downResult.ContainerStatus,
			ImageStatus:     downResult.ImageStatus,
		})
		if len(downs) == count {
			close(downResultCh)
		}
	}

	res := &api.DownResponse{
		Instances: downs,
	}

	return res, nil
}
