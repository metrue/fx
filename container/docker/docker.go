package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/apex/log"
	dockerTypes "github.com/docker/docker/api/types"
	dockerTypesContainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/metrue/fx/container"
	"github.com/metrue/fx/types"
)

// Docker manage container
type Docker struct {
	*client.Client
}

// CreateClient create a docker instance
func CreateClient() (*Docker, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	cli.NegotiateAPIVersion(ctx)
	return &Docker{cli}, nil
}

// Deploy create a Docker container from given image, and bind the constants.FxContainerExposePort to given port
func (d *Docker) Deploy(ctx context.Context, name string, image string, ports []int32) error {
	config := &dockerTypesContainer.Config{
		Image: image,
		ExposedPorts: nat.PortSet{
			"3000/tcp": struct{}{},
		},
	}

	bindings := []nat.PortBinding{}
	for _, port := range ports {
		bindings = append(bindings, nat.PortBinding{
			HostIP:   types.DefaultHost,
			HostPort: fmt.Sprintf("%d", port),
		})
	}

	hostConfig := &dockerTypesContainer.HostConfig{
		AutoRemove: true,
		PortBindings: nat.PortMap{
			"3000/tcp": bindings,
		},
	}

	reader, err := d.ImagePull(ctx, image, dockerTypes.ImagePullOptions{})
	if err != nil {
		return err
	}
	if os.Getenv("DEBUG") != "" {
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			return err
		}
		log.Info(string(body))
	}

	resp, err := d.ContainerCreate(ctx, config, hostConfig, nil, name)
	if os.Getenv("DEBUG") != "" {
		body, err := json.Marshal(resp)
		if err != nil {
			return err
		}
		log.Info(string(body))
	}

	if err := d.ContainerStart(ctx, resp.ID, dockerTypes.ContainerStartOptions{}); err != nil {
		return err
	}
	return nil
}

func (d *Docker) Update(ctx context.Context, name string) error {
	return nil
}

func (d *Docker) Destroy(ctx context.Context, name string) error {
	return d.ContainerStop(ctx, name, nil)
}

func (d *Docker) GetStatus(ctx context.Context, name string) error {
	return nil
}

var (
	_ container.Runner = &Docker{}
)
