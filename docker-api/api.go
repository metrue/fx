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
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type dockerInfo struct {
	Stream string `json:"stream"`
}

var cli *client.Client

func init() {
	var err error
	cli, err = client.NewEnvClient()
	if err != nil {
		panic(err)
	}
}

// Build builds a docker image from the image directory
func Build(name string, dir string) {
	tar := new(archivex.TarFile)
	tar.Create(dir)
	tar.AddAll(dir, false)
	tar.Close()
	dockerBuildContext, buildContextErr := os.Open(dir + ".tar")
	if buildContextErr != nil {
		panic(buildContextErr)
	}
	defer dockerBuildContext.Close()

	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile", // optional, is the default
		Tags:       []string{name},
		Labels:     map[string]string{"belong-to": "fx"},
	}
	buildResponse, buildErr := cli.ImageBuild(context.Background(), dockerBuildContext, buildOptions)
	if buildErr != nil {
		panic(buildErr)
	}
	log.Println("build", buildResponse.OSType)
	defer buildResponse.Body.Close()

	scanner := bufio.NewScanner(buildResponse.Body)
	for scanner.Scan() {
		var info dockerInfo
		err := json.Unmarshal(scanner.Bytes(), &info)
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Printf(info.Stream)
	}
}

// Pull image from hub.docker.com
func Pull(name string, verbose bool) {
	ctx := context.Background()
	if r, pullErr := cli.ImagePull(ctx, name, types.ImagePullOptions{}); pullErr != nil {
		panic(pullErr)
	} else {
		if verbose {
			io.Copy(os.Stdout, r)
		}
	}
}

// Deploy spins up a new container
func Deploy(name string, dir string, port string) {
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
		panic(err)
	}
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	fmt.Println(resp.ID)
}

// Stop interrupts a running container
func Stop(containerID string) (err error) {
	timeout := time.Duration(1) * time.Second
	err = cli.ContainerStop(context.Background(), containerID, &timeout)
	return err
}

// Remove interrupts and remove a running container
func Remove(containerID string) (err error) {
	return cli.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{Force: true})
}

// Remove remove docker image by imageID
func ImageRemove(imageID string) (err error) {
	_, err = cli.ImageRemove(context.Background(), imageID, types.ImageRemoveOptions{Force: true})
	return err
}
