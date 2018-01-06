package docker

import (
	"bufio"
	"encoding/json"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/jhoonb/archivex"
	"github.com/pkg/errors"

	"context"
	"fmt"
	"io"
	"os"
	"time"
)

type dockerInfo struct {
	Stream string `json:"stream"`
}

var docekClient *client.Client

func getClient() (*client.Client, error) {
	if docekClient != nil {
		return docekClient, nil
	}
	var err error
	docekClient, err = client.NewEnvClient()
	return docekClient, err
}

func Info() error {
	cli, err := getClient()
	if err != nil {
		return errors.Wrap(err, "Create Docker client failed")
	}
	ctx := context.Background()
	_, err = cli.Info(ctx)
	return err
}

// Build builds a docker image from the image directory
func Build(name string, dir string) error {
	cli, err := getClient()
	if err != nil {
		return err
	}

	tar := new(archivex.TarFile)
	err = tar.Create(dir)
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
	buildResponse, buildErr := cli.ImageBuild(context.Background(), dockerBuildContext, buildOptions)
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
		// fmt.Printf(info.Stream)
	}

	return nil
}

// Pull image from hub.docker.com
func Pull(name string, verbose bool) error {
	cli, err := getClient()
	if err != nil {
		return err
	}

	ctx := context.Background()
	r, pullErr := cli.ImagePull(ctx, name, types.ImagePullOptions{})
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
	cli, err := getClient()
	if err != nil {
		return nil, err
	}

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

	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, "")
	if err != nil {
		return nil, err
	}

	if err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}

	fmt.Printf("Deployed to container %s\n", resp.ID)
	return &resp, err
}

// Stop interrupts a running container
func Stop(containerID string) (err error) {
	cli, err := getClient()
	if err != nil {
		return err
	}
	timeout := time.Duration(1) * time.Second
	err = cli.ContainerStop(context.Background(), containerID, &timeout)
	return err
}

// Remove interrupts and remove a running container
func Remove(containerID string) (err error) {
	cli, err := getClient()
	if err != nil {
		return err
	}
	return cli.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{Force: true})
}

//ImageRemove remove docker image by imageID
func ImageRemove(imageID string) (err error) {
	cli, err := getClient()
	if err != nil {
		return err
	}
	_, err = cli.ImageRemove(context.Background(), imageID, types.ImageRemoveOptions{Force: true})
	return err
}
