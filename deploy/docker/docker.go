package docker

import (
	"context"
	"os"

	"github.com/metrue/fx/constants"
	containerruntimes "github.com/metrue/fx/container_runtimes"
	dockerHTTP "github.com/metrue/fx/container_runtimes/docker/http"
	dockerSDK "github.com/metrue/fx/container_runtimes/docker/sdk"
	"github.com/metrue/fx/deploy"
	"github.com/metrue/fx/types"
)

// Docker manage container
type Docker struct {
	cli containerruntimes.ContainerRuntime
}

// CreateClient create a docker instance
func CreateClient(ctx context.Context) (d *Docker, err error) {
	var cli containerruntimes.ContainerRuntime
	host := os.Getenv("DOCKER_REMOTE_HOST_ADDR")
	user := os.Getenv("DOCKER_REMOTE_HOST_USER")
	if host != "" && user != "" {
		cli, err = dockerHTTP.Create(host, constants.AgentPort)
		if err != nil {
			return nil, err
		}
	} else {
		cli, err = dockerSDK.CreateClient(ctx)
		if err != nil {
			return nil, err
		}
	}

	return &Docker{cli: cli}, nil
}

// Deploy create a Docker container from given image, and bind the constants.FxContainerExposePort to given port
func (d *Docker) Deploy(ctx context.Context, fn types.Func, name string, image string, ports []types.PortBinding) error {
	return d.cli.StartContainer(ctx, name, image, ports)
}

// Update a container
func (d *Docker) Update(ctx context.Context, name string) error {
	return nil
}

// Destroy stop and remove container
func (d *Docker) Destroy(ctx context.Context, name string) error {
	return d.cli.StopContainer(ctx, name)
}

// GetStatus get status of container
func (d *Docker) GetStatus(ctx context.Context, name string) error {
	return nil
}

// List services
func (d *Docker) List(ctx context.Context, name string) ([]types.Service, error) {
	// FIXME support remote host
	return d.cli.ListContainer(ctx, name)
}

var (
	_ deploy.Deployer = &Docker{}
)
