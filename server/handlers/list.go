package handlers

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// List returns the list of running services
func List() []types.Container {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	return containers
}
