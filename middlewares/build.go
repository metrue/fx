package middlewares

import (
	"fmt"
	"os"
	"time"

	"github.com/metrue/fx/bundle"
	"github.com/metrue/fx/constants"
	dockerHTTP "github.com/metrue/fx/container_runtimes/docker/http"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/hook"
	"github.com/metrue/fx/packer"
	"github.com/metrue/fx/pkg/spinner"
	"github.com/metrue/fx/utils"
)

// Build image
func Build(ctx context.Contexter) (err error) {
	const task = "building"
	spinner.Start(task)
	defer func() {
		spinner.Stop(task, err)
	}()

	workdir := fmt.Sprintf("/tmp/fx-%d", time.Now().Unix())
	if err := utils.EnsureDir(workdir); err != nil {
		return err
	}
	defer os.RemoveAll(workdir)

	// Cases supports
	// 1. a single file function
	// 		fx up func.js
	// 2. a directory with Docker in it
	// 		fx up ./func/
	// 3. a directory without Dockerfile in it, but has fx handle function file
	// 4. a fx handlefunction file and its dependencies files or/and directory
	// 		fx up func.js helper.js ./lib/

	// When only one directory given and there is a Dockerfile in given directory, treat it as a containerized project and skip packing
	fn := ctx.Get("fn").(string)
	deps := ctx.Get("deps").([]string)
	language := ctx.Get("language").(string)
	host := ctx.Get("host").(string)
	kubeconf := ctx.Get("kubeconf").(string)
	name := fmt.Sprintf("%s-%s", ctx.Get("name").(string), time.Now().Format("20060102150405"))

	if err := bundle.Bundle(workdir, language, fn, deps...); err != nil {
		return err
	}

	if err := hook.RunBeforeBuildHook(workdir); err != nil {
		return err
	}

	if host != "" {
		// TODO port should be configurable
		docker, err := dockerHTTP.Create(host, constants.AgentPort)
		if err != nil {
			return err
		}
		if err := docker.BuildImage(ctx.GetContext(), workdir, name); err != nil {
			return err
		}
		nameWithTag := fmt.Sprintf("%s:latest", name)
		if err := docker.TagImage(ctx.GetContext(), name, nameWithTag); err != nil {
			return err
		}
		ctx.Set("image", nameWithTag)
	}

	if kubeconf != "" {
		data, err := packer.PackIntoK8SConfigMapFile(workdir)
		if err != nil {
			return err
		}
		ctx.Set("data", data)
	}

	return nil
}
