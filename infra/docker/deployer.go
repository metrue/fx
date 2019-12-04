package docker

import (
	"context"
	"strconv"

	dockerTypes "github.com/docker/docker/api/types"
	containerruntimes "github.com/metrue/fx/container_runtimes"
	"github.com/metrue/fx/infra"
	"github.com/metrue/fx/pkg/spinner"
	"github.com/metrue/fx/types"
)

// Deployer manage container
type Deployer struct {
	cli containerruntimes.ContainerRuntime
}

// CreateClient create a docker instance
func CreateClient(client containerruntimes.ContainerRuntime) (d *Deployer, err error) {
	return &Deployer{cli: client}, nil
}

// Deploy create a Docker container from given image, and bind the constants.FxContainerExposePort to given port
func (d *Deployer) Deploy(ctx context.Context, fn types.Func, name string, image string, ports []types.PortBinding) (err error) {
	spinner.Start("deploying " + name)
	defer func() {
		spinner.Stop("deploying "+name, err)
	}()
	return d.cli.StartContainer(ctx, name, image, ports)
}

// Update a container
func (d *Deployer) Update(ctx context.Context, name string) error {
	return nil
}

// Destroy stop and remove container
func (d *Deployer) Destroy(ctx context.Context, name string) (err error) {
	spinner.Start("destroying " + name)
	defer func() {
		spinner.Stop("destroying "+name, err)
	}()
	return d.cli.StopContainer(ctx, name)
}

// GetStatus get a service status
func (d *Deployer) GetStatus(ctx context.Context, name string) (types.Service, error) {
	var container dockerTypes.ContainerJSON
	if err := d.cli.InspectContainer(ctx, name, &container); err != nil {
		return types.Service{}, err
	}

	service := types.Service{
		ID:   container.ID,
		Name: container.Name,
	}
	for _, bindings := range container.NetworkSettings.Ports {
		if len(bindings) > 0 {
			binding := bindings[0]
			port, err := strconv.Atoi(binding.HostPort)
			if err != nil {
				return service, err
			}
			service.Port = port
			service.Host = binding.HostIP
			service.State = container.State.Status
			service.Image = container.Image
			break
		}
		if service.Port != 0 && service.Host != "" {
			break
		}
	}

	return service, nil
}

// Ping check healty status of infra
func (d *Deployer) Ping(ctx context.Context) error {
	if _, err := d.cli.Version(ctx); err != nil {
		return err
	}
	return nil
}

// List services
func (d *Deployer) List(ctx context.Context, name string) (svcs []types.Service, err error) {
	const task = "listing"
	spinner.Start(task)
	defer func() {
		spinner.Stop(task, err)
	}()

	// FIXME support remote host
	return d.cli.ListContainer(ctx, name)
}

var (
	_ infra.Deployer = &Deployer{}
)
