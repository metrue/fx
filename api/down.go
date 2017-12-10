package api

import (
	"github.com/metrue/fx/handlers"
	Message "github.com/metrue/fx/message"
)

func Down(req *DownRequest) (*DownResponse, error) {

	containers := handlers.List(req.ID...)
	count := len(containers)
	downResultCh := make(chan Message.DownMsgMeta, count)

	for _, container := range containers {
		go handlers.Down(container.ID[:10], container.Image, downResultCh)
	}

	// collect down result
	var downs []*DownMsgMeta
	for downResult := range downResultCh {
		//TODO use only one type avoiding conversion where possible
		downs = append(downs, &DownMsgMeta{
			ContainerId:     downResult.ContainerId,
			ContainerStatus: downResult.ContainerStatus,
			ImageStatus:     downResult.ImageStatus,
		})
		if len(downs) == count {
			close(downResultCh)
		}
	}

	res := &DownResponse{
		Instances: downs,
	}

	return res, nil
}
