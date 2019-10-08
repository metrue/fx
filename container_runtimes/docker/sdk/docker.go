package docker

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/apex/log"
	dockerTypes "github.com/docker/docker/api/types"
	dockerTypesContainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	containerruntimes "github.com/metrue/fx/container_runtimes"
	"github.com/metrue/fx/types"
	"github.com/metrue/fx/utils"
)

// Docker docker as image builder
type Docker struct {
	*client.Client
}

// CreateClient create a docker instance
func CreateClient(ctx context.Context) (*Docker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	cli.NegotiateAPIVersion(ctx)
	return &Docker{cli}, nil
}

// BuildImage a directory to be a image
func (d *Docker) BuildImage(ctx context.Context, workdir string, name string) error {
	tarDir, err := ioutil.TempDir("/tmp", "fx-tar")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tarDir)

	imageID := uuid.New().String()
	tarFilePath := filepath.Join(tarDir, fmt.Sprintf("%s.tar", imageID))

	if err := utils.TarDir(workdir, tarFilePath); err != nil {
		return err
	}

	dockerBuildContext, err := os.Open(tarFilePath)
	if err != nil {
		return err
	}
	defer dockerBuildContext.Close()

	options := dockerTypes.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{imageID, name},
		Labels: map[string]string{
			"belong-to": "fx",
		},
	}

	resp, err := d.ImageBuild(ctx, dockerBuildContext, options)
	if err != nil {
		return err
	}

	if os.Getenv("DEBUG") != "" {
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		log.Info(string(body))
	}

	return nil
}

// PushImage push image to hub.docker.com
func (d *Docker) PushImage(ctx context.Context, name string) (string, error) {
	username := os.Getenv("DOCKER_USERNAME")
	password := os.Getenv("DOCKER_PASSWORD")
	if username == "" || password == "" {
		return "", fmt.Errorf("DOCKER_USERNAME and DOCKER_PASSWORD required for push image to registy")
	}

	// TODO support private registy, like Azure Container registry
	authConfig := dockerTypes.AuthConfig{
		Username: username,
		Password: password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}

	nameWithTag := username + "/" + name
	if err := d.ImageTag(ctx, name, nameWithTag); err != nil {
		return "", err
	}

	options := dockerTypes.ImagePushOptions{
		RegistryAuth: base64.URLEncoding.EncodeToString(encodedJSON),
	}
	resp, err := d.ImagePush(ctx, nameWithTag, options)
	if err != nil {
		return "", err
	}
	defer resp.Close()

	if os.Getenv("DEBUG") != "" {
		body, err := ioutil.ReadAll(resp)
		if err != nil {
			return "", err
		}
		log.Info(string(body))
	}
	return nameWithTag, nil
}

// InspectImage inspect a image
func (d *Docker) InspectImage(ctx context.Context, name string, img interface{}) error {
	_, body, err := d.ImageInspectWithRaw(ctx, name)
	if err != nil {
		return err
	}
	rdr := bytes.NewReader(body)
	return json.NewDecoder(rdr).Decode(&img)
}

// StartContainer create and start a container from given image
func (d *Docker) StartContainer(ctx context.Context, name string, image string, ports []int32) error {
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

// StopContainer stop and remove container
func (d *Docker) StopContainer(ctx context.Context, name string) error {
	return d.ContainerStop(ctx, name, nil)
}

// InspectContainer inspect a container
func (d *Docker) InspectContainer(ctx context.Context, name string, container interface{}) error {
	return nil
}

var (
	_ containerruntimes.ContainerRuntime = &Docker{}
)
