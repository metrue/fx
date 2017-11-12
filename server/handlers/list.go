package handlers

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

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
