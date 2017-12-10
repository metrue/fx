package api

import (
	"fmt"

	"github.com/metrue/fx/handlers"
)

func List(req *ListRequest) (*ListResponse, error) {

	containers := handlers.List(req.ID...)
	var list []*ListItem
	for _, container := range containers {

		var serviceURL string
		if len(container.Ports) > 0 {
			serviceURL = fmt.Sprintf("%s:%d", container.Ports[0].IP, container.Ports[0].PublicPort)
		}

		item := &ListItem{
			FunctionID: container.ID[:10],
			ServiceURL: serviceURL,
			State:      container.State,
		}

		list = append(list, item)
	}

	res := &ListResponse{
		Instances: list,
	}

	return res, nil
}
