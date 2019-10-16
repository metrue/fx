package docker

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	dockerTypes "github.com/docker/docker/api/types"
	runtime "github.com/metrue/fx/container_runtimes/docker/sdk"
	"github.com/metrue/fx/deploy"
	"github.com/metrue/fx/packer"
	"github.com/metrue/fx/types"
	"github.com/metrue/fx/utils"
)

// Docker manage container
type Docker struct {
	client *runtime.Docker
}

// CreateClient create a docker instance
func CreateClient(ctx context.Context) (*Docker, error) {
	cli, err := runtime.CreateClient(ctx)
	if err != nil {
		return nil, err
	}
	return &Docker{client: cli}, nil
}

// Deploy create a Docker container from given image, and bind the constants.FxContainerExposePort to given port
func (d *Docker) Deploy(ctx context.Context, fn types.Func, name string, ports []types.PortBinding) error {
	workdir := fmt.Sprintf("/tmp/fx-%d", time.Now().Unix())
	defer os.RemoveAll(workdir)

	if err := packer.PackIntoDir(fn, workdir); err != nil {
		return err
	}

	if err := d.client.BuildImage(ctx, workdir, name); err != nil {
		return err
	}

	nameWithTag := name + ":latest"
	if err := d.client.ImageTag(ctx, name, nameWithTag); err != nil {
		log.Fatalf("could tag image: %v", err)
		return err
	}

	// when deploy a function on a bare Docker running without Kubernetes,
	// image would be built on-demand on host locally, so there is no need to
	// pull image from remote.
	// But it takes some times waiting image ready after image built, we retry to make sure it ready here
	var imgInfo dockerTypes.ImageInspect
	if err := utils.RunWithRetry(func() error {
		return d.client.InspectImage(ctx, name, &imgInfo)
	}, time.Second*1, 5); err != nil {
		return err
	}

	return d.client.StartContainer(ctx, name, name, ports)
}

// Update a container
func (d *Docker) Update(ctx context.Context, name string) error {
	return nil
}

// Destroy stop and remove container
func (d *Docker) Destroy(ctx context.Context, name string) error {
	return d.client.ContainerStop(ctx, name, nil)
}

// GetStatus get status of container
func (d *Docker) GetStatus(ctx context.Context, name string) error {
	return nil
}

var (
	_ deploy.Deployer = &Docker{}
)
