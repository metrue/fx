package docker

import (
	"context"

	containerruntimes "github.com/metrue/fx/container_runtimes"
	"github.com/metrue/fx/deploy"
	"github.com/metrue/fx/types"
)

// Docker manage container
type Docker struct {
	cli containerruntimes.ContainerRuntime
}

// CreateClient create a docker instance
func CreateClient(client containerruntimes.ContainerRuntime) (d *Docker, err error) {
	return &Docker{cli: client}, nil
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
