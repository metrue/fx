package middlewares

import (
	"fmt"
	"os"
	"time"

	containerruntimes "github.com/metrue/fx/container_runtimes"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/packer"
	"github.com/metrue/fx/pkg/spinner"
)

// Build image
func Build(ctx context.Contexter) (err error) {
	const task = "building"
	spinner.Start(task)
	defer func() {
		spinner.Stop(task, err)
	}()

	name := ctx.Get("name").(string)
	docker := ctx.Get("docker").(containerruntimes.ContainerRuntime)
	workdir := fmt.Sprintf("/tmp/fx-%d", time.Now().Unix())
	defer os.RemoveAll(workdir)

	if err := packer.Pack(workdir, ctx.Get("sources").([]string)...); err != nil {
		return err
	}

	data, err := packer.PackIntoK8SConfigMapFile(workdir)
	if err != nil {
		return err
	}
	ctx.Set("data", data)

	if err := docker.BuildImage(ctx.GetContext(), workdir, name); err != nil {
		return err
	}

	nameWithTag := name + ":latest"
	if err := docker.TagImage(ctx.GetContext(), name, nameWithTag); err != nil {
		return err
	}
	ctx.Set("image", nameWithTag)

	if os.Getenv("K3S") != "" {
		username := os.Getenv("DOCKER_USERNAME")
		password := os.Getenv("DOCKER_PASSWORD")
		if username != "" && password != "" {
			if _, err := docker.PushImage(ctx.GetContext(), name); err != nil {
				return err
			}
			ctx.Set("image", username+"/"+name)
		}
	}
	return nil
}
