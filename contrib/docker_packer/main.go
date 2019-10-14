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

// Version binary version
var Version = "0.0.1"

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
		os.Exit(1)
	}

	var tree map[string]string
	if err := json.Unmarshal([]byte(str), &tree); err != nil {
		log.Fatalf("could not unmarshal meta: %s", meta)
		os.Exit(1)
	}
	workdir := "/tmp/fx"
	if err := packer.TreeToDir(tree, workdir); err != nil {
		log.Fatalf("could not restore to dir: %v", err)
		os.Exit(1)
	}
	defer os.RemoveAll(workdir)

	ctx := context.Background()
	dockerClient, err := runtime.CreateClient(ctx)
	if err != nil {
		log.Fatalf("could not create a docker client: %v", err)
		os.Exit(1)
	}
	if err := dockerClient.BuildImage(ctx, workdir, name); err != nil {
		log.Fatalf("could not build image: %s", err)
		os.Exit(1)
	}

	nameWithTag := name + ":latest"
	if err := dockerClient.ImageTag(ctx, name, nameWithTag); err != nil {
		log.Fatalf("could tag image: %v", err)
		os.Exit(1)
	}
	var imgInfo dockerTypes.ImageInspect
	if err := utils.RunWithRetry(func() error {
		return dockerClient.InspectImage(context.Background(), name, &imgInfo)
	}, time.Second*1, 5); err != nil {
		fmt.Printf("inspect image failed: %s", err)
	}
	fmt.Println("image built succcessfully")
}
