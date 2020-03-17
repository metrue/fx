package middlewares

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/metrue/fx/bundle"
	"github.com/metrue/fx/constants"
	dockerHTTP "github.com/metrue/fx/container_runtimes/docker/http"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/hook"
	dockerInfra "github.com/metrue/fx/infra/docker"
	k8sInfra "github.com/metrue/fx/infra/k8s"
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
	name := ctx.Get("name").(string)

	if err := bundle.Bundle(workdir, language, fn, deps...); err != nil {
		return err
	}

	if err := hook.RunBeforeBuildHook(workdir); err != nil {
		return err
	}

	if host != "" {
		addr := strings.Split(host, "@")
		if len(addr) != 2 {
			return fmt.Errorf("invalid host information, should be format of <user>@<ip>")
		}
		ip := addr[1]
		// TODO port should be configurable
		docker, err := dockerHTTP.Create(ip, constants.AgentPort)
		if err != nil {
			return err
		}
		if err := docker.BuildImage(ctx.GetContext(), workdir, name); err != nil {
			return err
		}
		nameWithTag := name + ":latest"
		if err := docker.TagImage(ctx.GetContext(), name, nameWithTag); err != nil {
			return err
		}
		ctx.Set("image", nameWithTag)

		deployer, err := dockerInfra.CreateDeployer(docker)
		if err != nil {
			return err
		}
		ctx.Set("docker_deployer", deployer)
	}

	if kubeconf != "" {
		data, err := packer.PackIntoK8SConfigMapFile(workdir)
		if err != nil {
			return err
		}
		ctx.Set("data", data)

		deployer, err := k8sInfra.CreateDeployer(kubeconf)
		if err != nil {
			return err
		}
		ctx.Set("k8s_deployer", deployer)
	}

	return nil
}
