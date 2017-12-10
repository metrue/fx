package service

import (
	"fmt"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/handlers"
)

func List(req *api.ListRequest) (*api.ListResponse, error) {

	containers := handlers.List(req.ID...)
	var list []*api.ListItem
	for _, container := range containers {

		var serviceURL string
		if len(container.Ports) > 0 {
			serviceURL = fmt.Sprintf("%s:%d", container.Ports[0].IP, container.Ports[0].PublicPort)
		}

		item := &api.ListItem{
			FunctionID: container.ID[:10],
			ServiceURL: serviceURL,
			State:      container.State,
		}

		list = append(list, item)
	}

	res := &api.ListResponse{
		Instances: list,
	}

	return res, nil
}
