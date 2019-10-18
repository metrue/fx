package docker

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	dockerTypes "github.com/docker/docker/api/types"
	dockerHTTP "github.com/metrue/fx/container_runtimes/docker/http"
	runtime "github.com/metrue/fx/container_runtimes/docker/sdk"
	"github.com/metrue/fx/deploy"
	"github.com/metrue/fx/packer"
	"github.com/metrue/fx/types"
	"github.com/metrue/fx/utils"
	"github.com/pkg/errors"
)

// Docker manage container
type Docker struct {
	localClient *runtime.Docker
}

// CreateClient create a docker instance
func CreateClient(ctx context.Context) (*Docker, error) {
	cli, err := runtime.CreateClient(ctx)
	if err != nil {
		return nil, err
	}
	return &Docker{localClient: cli}, nil
}

// Deploy create a Docker container from given image, and bind the constants.FxContainerExposePort to given port
func (d *Docker) Deploy(ctx context.Context, fn types.Func, name string, ports []types.PortBinding) error {
	// if DOCKER_REMOTE_HOST and DOCKER_REMOTE_PORT given
	// it means user is going to deploy service to remote host
	host := os.Getenv("DOCKER_REMOTE_HOST")
	port := os.Getenv("DOCKER_REMOTE_PORT")
	if port != "" && host != "" {
		httpClient, err := dockerHTTP.Create(host, port)
		if err != nil {
			return err
		}

		project, err := packer.Pack(name, fn)
		if err != nil {
			return errors.Wrapf(err, "could pack function %v (%s)", name, fn)
		}
		return httpClient.Up(dockerHTTP.UpOptions{
			Body:       []byte(fn.Source),
			Lang:       fn.Language,
			Name:       name,
			Port:       int(ports[0].ServiceBindingPort),
			HealtCheck: false,
			Project:    project,
		})
	}

	workdir := fmt.Sprintf("/tmp/fx-%d", time.Now().Unix())
	defer os.RemoveAll(workdir)

	if err := packer.PackIntoDir(fn, workdir); err != nil {
		log.Fatalf("could not pack function %v: %v", fn, err)
		return err
	}
	if err := d.localClient.BuildImage(ctx, workdir, name); err != nil {
		log.Fatalf("could not build image: %v", err)
		return err
	}

	nameWithTag := name + ":latest"
	if err := d.localClient.ImageTag(ctx, name, nameWithTag); err != nil {
		log.Fatalf("could not tag image: %v", err)
		return err
	}

	// when deploy a function on a bare Docker running without Kubernetes,
	// image would be built on-demand on host locally, so there is no need to
	// pull image from remote.
	// But it takes some times waiting image ready after image built, we retry to make sure it ready here
	var imgInfo dockerTypes.ImageInspect
	if err := utils.RunWithRetry(func() error {
		return d.localClient.InspectImage(ctx, name, &imgInfo)
	}, time.Second*1, 5); err != nil {
		return err
	}

	return d.localClient.StartContainer(ctx, name, name, ports)
}

// Update a container
func (d *Docker) Update(ctx context.Context, name string) error {
	return nil
}

// Destroy stop and remove container
func (d *Docker) Destroy(ctx context.Context, name string) error {
	return d.localClient.ContainerStop(ctx, name, nil)
}

// GetStatus get status of container
func (d *Docker) GetStatus(ctx context.Context, name string) error {
	return nil
}

var (
	_ deploy.Deployer = &Docker{}
)
