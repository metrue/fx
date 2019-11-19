package middlewares

import (
	"fmt"
	"os"
	"time"

	containerruntimes "github.com/metrue/fx/container_runtimes"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/packer"
	"github.com/metrue/fx/pkg/spinner"
	"github.com/metrue/fx/types"
)

// Build image
func Build(ctx *context.Context) (err error) {
	const task = "building"
	spinner.Start(task)
	defer func() {
		spinner.Stop(task, err)
	}()

	cli := ctx.GetCliContext()
	name := cli.String("name")
	fn := ctx.Get("fn").(types.Func)
	docker := ctx.Get("docker").(containerruntimes.ContainerRuntime)
	workdir := fmt.Sprintf("/tmp/fx-%d", time.Now().Unix())
	defer os.RemoveAll(workdir)

	if err := packer.PackIntoDir(fn, workdir); err != nil {
		return err
	}
	if err := docker.BuildImage(ctx.Context, workdir, name); err != nil {
		return err
	}

	nameWithTag := name + ":latest"
	if err := docker.TagImage(ctx, name, nameWithTag); err != nil {
		return err
	}
	ctx.Set("image", nameWithTag)

	if os.Getenv("K3S") != "" {
		name := cli.String("name")
		username := os.Getenv("DOCKER_USERNAME")
		password := os.Getenv("DOCKER_PASSWORD")
		if username != "" && password != "" {
			if _, err := docker.PushImage(ctx.Context, name); err != nil {
				return err
			}
			ctx.Set("image", username+"/"+name)
		}
	}
	return nil
}
