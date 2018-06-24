package service

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/metrue/fx/api"
	"github.com/pkg/errors"
)

// List returns the list of running services
func DoList(containerIds ...string) ([]types.Container, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		err = errors.Wrap(err, "[list.go] New client failed")
		return nil, err
	}

	filters := filters.NewArgs()
	filters.Add("label", "belong-to=fx")
	filters.Add("status", "running")
	if len(containerIds) > 0 {
		for _, id := range containerIds {
			filters.Add("id", id)
		}
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{
		All:     true,
		Filters: filters,
	})
	if err != nil {
		err = errors.Wrap(err, "[list.go] list containers failed")
		return nil, err
	}

	return containers, nil
}

//List handles a functions listing request
func List(ctx context.Context, req *api.ListRequest) (*api.ListResponse, error) {

	containers, err := DoList(req.ID...)
	if err != nil {
		return nil, err
	}
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
