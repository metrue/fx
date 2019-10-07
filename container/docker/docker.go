package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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

	// when deploy a function on a bare Docker running without Kubernetes,
	// image would be built on-demand on host locally, so there is no need to
	// pull image from remote.
	if _, ok := d.isImageExisted(ctx, image); !ok {
		fmt.Println("++++++++++")
		fmt.Println(image, " ---> not ready")
		fmt.Println("++++++++++")
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
	}

	resp, err := d.ContainerCreate(ctx, config, hostConfig, nil, name)
	if os.Getenv("DEBUG") != "" {
		body, err := json.Marshal(resp)
		if err != nil {
			return err
		}
		log.Info(string(body))
	}
	if err != nil {
		return err
	}

	if err := d.ContainerStart(ctx, resp.ID, dockerTypes.ContainerStartOptions{}); err != nil {
		return err
	}
	return nil
}

// Update a container
func (d *Docker) Update(ctx context.Context, name string) error {
	// TODO
	// UpdateConfig holds the mutable attributes of a Container.
	// Those attributes can be updated at runtime.
	updateConfig := dockerTypesContainer.UpdateConfig{}
	resp, err := d.ContainerUpdate(ctx, name, updateConfig)
	if os.Getenv("DEBUG") != "" {
		log.Infof("%s", resp.Warnings)
	}
	if err != nil {
		return err
	}
	return nil
}

// Destroy stop and remove container
func (d *Docker) Destroy(ctx context.Context, name string) error {
	return d.ContainerStop(ctx, name, nil)
}

// GetStatus get status of container
func (d *Docker) GetStatus(ctx context.Context, name string) error {
	return nil
}

func (d *Docker) isImageExisted(ctx context.Context, name string) (string, bool) {
	images, err := d.ImageList(ctx, dockerTypes.ImageListOptions{})
	if err != nil {
		log.Warnf("list images failed: %v", err)
		return "", false
	}

	var tag string
	var found bool
	for _, img := range images {
		for _, fullTag := range img.RepoTags {
			fmt.Println("--->", fullTag)
			arr := strings.Split(fullTag, ":")
			if arr[0] == name {
				found = true
				tag = fullTag
				break
			}
		}
	}
	return tag, found
}

var (
	_ container.Runner = &Docker{}
)
