package docker

import (
	"bufio"
	"encoding/json"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/jhoonb/archivex"

	"context"
	"io"
	"os"
	"time"
)

type dockerInfo struct {
	Stream string `json:"stream"`
}

var dockerClient *client.Client

func init() {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	dockerClient = cli
}

func IsRunning() bool {
	ctx := context.Background()
	_, err := dockerClient.Info(ctx)
	return err == nil
}

// Build builds a docker image from the image directory
func Build(name string, dir string) error {
	tar := new(archivex.TarFile)
	err := tar.Create(dir)
	if err != nil {
		return err
	}
	err = tar.AddAll(dir, false)
	if err != nil {
		return err
	}
	err = tar.Close()
	if err != nil {
		return err
	}

	dockerBuildContext, buildContextErr := os.Open(dir + ".tar")
	if buildContextErr != nil {
		return buildContextErr
	}
	defer dockerBuildContext.Close()

	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile", // optional, is the default
		Tags:       []string{name},
		Labels:     map[string]string{"belong-to": "fx"},
	}
	buildResponse, buildErr := dockerClient.ImageBuild(context.Background(), dockerBuildContext, buildOptions)
	if buildErr != nil {
		return buildErr
	}
	defer buildResponse.Body.Close()

	scanner := bufio.NewScanner(buildResponse.Body)
	for scanner.Scan() {
		var info dockerInfo
		err := json.Unmarshal(scanner.Bytes(), &info)
		if err != nil {
			return err
		}
	}

	return nil
}

// Pull image from hub.docker.com
func Pull(name string, verbose bool) error {
	ctx := context.Background()
	r, pullErr := dockerClient.ImagePull(ctx, name, types.ImagePullOptions{})
	if pullErr != nil {
		return pullErr
	}

	if verbose {
		io.Copy(os.Stdout, r)
	}

	return nil
}

// Deploy spins up a new container
func Deploy(name string, dir string, port string) (*container.ContainerCreateCreatedBody, error) {
	ctx := context.Background()
	imageName := name
	containerConfig := &container.Config{
		Image: imageName,
		ExposedPorts: nat.PortSet{
			"3000/tcp": struct{}{},
		},
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"3000/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: port,
				},
			},
		},
	}

	resp, err := dockerClient.ContainerCreate(ctx, containerConfig, hostConfig, nil, "")
	if err != nil {
		return nil, err
	}

	if err = dockerClient.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	return &resp, err
}

// Stop interrupts a running container
func Stop(containerID string) (err error) {
	timeout := time.Duration(1) * time.Second
	err = dockerClient.ContainerStop(context.Background(), containerID, &timeout)
	return err
}

// Remove interrupts and remove a running container
func Remove(containerID string) (err error) {
	return dockerClient.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{Force: true})
}

//ImageRemove remove docker image by imageID
func ImageRemove(imageID string) (err error) {
	_, err = dockerClient.ImageRemove(context.Background(), imageID, types.ImageRemoveOptions{Force: true})
	return err
}
