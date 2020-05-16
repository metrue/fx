package docker

import (
	"context"
	"strconv"

	dockerTypes "github.com/docker/docker/api/types"
	containerruntimes "github.com/metrue/fx/container_runtimes"
	"github.com/metrue/fx/driver"
	"github.com/metrue/fx/pkg/spinner"
	"github.com/metrue/fx/types"
)

// Driver manage container
type Driver struct {
	dockerClient containerruntimes.ContainerRuntime
}

// Options to initialize a fx docker driver
type Options struct {
	DockerClient containerruntimes.ContainerRuntime
}

// New a fx docker driver
func New(options Options) *Driver {
	return &Driver{
		dockerClient: options.DockerClient,
	}
}

// Ping check healty status of driver
func (d *Driver) Ping(ctx context.Context) error {
	if _, err := d.dockerClient.Version(ctx); err != nil {
		return err
	}
	return nil
}

// Deploy create a Docker container from given image, and bind the constants.FxContainerExposePort to given port
func (d *Driver) Deploy(ctx context.Context, fn string, name string, image string, ports []types.PortBinding) (err error) {
	spinner.Start("deploying " + name)
	defer func() {
		spinner.Stop("deploying "+name, err)
	}()
	return d.dockerClient.StartContainer(ctx, name, image, ports)
}

// Update a container
func (d *Driver) Update(ctx context.Context, name string) error {
	return nil
}

// Destroy stop and remove container
func (d *Driver) Destroy(ctx context.Context, name string) (err error) {
	spinner.Start("destroying " + name)
	defer func() {
		spinner.Stop("destroying "+name, err)
	}()
	if err := d.dockerClient.StopContainer(ctx, name); err != nil {
		return err
	}
	return d.dockerClient.RemoveContainer(ctx, name)
}

// GetStatus get a service status
func (d *Driver) GetStatus(ctx context.Context, name string) (types.Service, error) {
	var container dockerTypes.ContainerJSON
	if err := d.dockerClient.InspectContainer(ctx, name, &container); err != nil {
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

// List services
func (d *Driver) List(ctx context.Context, name string) (svcs []types.Service, err error) {
	// FIXME support remote host
	return d.dockerClient.ListContainer(ctx, name)
}

var (
	_ driver.Driver = &Driver{}
)
