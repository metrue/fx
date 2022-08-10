package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	dockerTypes "github.com/docker/docker/api/types"
	runtime "github.com/metrue/fx/container_runtimes/docker/sdk"
	"github.com/metrue/fx/packer"
	"github.com/metrue/fx/utils"
)

func init() {
	// TODO clean it up
	os.Setenv("DEBUG", "true")
}

func main() {
	args := os.Args

	if len(args) != 3 {
		fmt.Println(`Usage:
docker_packer <encrypt_docker_project_source_tree> <image_name>
		`)
		return
	}

	meta := args[1]
	name := args[2]

	str, err := base64.StdEncoding.WithPadding(base64.StdPadding).DecodeString(meta)
	if err != nil {
		log.Fatalf("could decode meta: %s, %v", meta, err)
	}

	var tree map[string]string
	//nolint
	if err := json.Unmarshal([]byte(str), &tree); err != nil {
		log.Fatalf("could not unmarshal meta: %s", meta)
	}
	workdir := "/tmp/fx"
	if err := packer.TreeToDir(tree, workdir); err != nil {
		log.Fatalf("could not restore to dir: %v", err)
	}
	defer os.RemoveAll(workdir)

	ctx := context.Background()
	dockerClient, err := runtime.CreateClient(ctx)
	if err != nil {
		log.Fatalf("could not create a docker client: %v", err)
	}
	if err := dockerClient.BuildImage(ctx, workdir, name); err != nil {
		log.Fatalf("could not build image: %s", err)
	}

	nameWithTag := name + ":latest"
	if err := dockerClient.ImageTag(ctx, name, nameWithTag); err != nil {
		log.Fatalf("could tag image: %v", err)
	}
	var imgInfo dockerTypes.ImageInspect
	if err := utils.RunWithRetry(func() error {
		return dockerClient.InspectImage(context.Background(), name, &imgInfo)
	}, time.Second*1, 5); err != nil {
		fmt.Printf("inspect image failed: %s", err)
	}
	fmt.Println("image built succcessfully")
}
