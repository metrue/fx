package handlers

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// List returns the list of running services
func List(containerIds ...string) ([]types.Container, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
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
		return nil, err
	}
	return containers, nil
}
