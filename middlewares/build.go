package middlewares

import (
	"fmt"
	"os"
	"time"

	containerruntimes "github.com/metrue/fx/container_runtimes"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/hook"
	"github.com/metrue/fx/packer"
	"github.com/metrue/fx/pkg/spinner"
	"github.com/metrue/fx/types"
	"github.com/metrue/fx/utils"
	"github.com/otiai10/copy"
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
	sources := ctx.Get("sources").([]string)

	if len(sources) == 0 {
		return fmt.Errorf("source file/directory of function required")
	}

	if len(sources) == 1 &&
		utils.IsDir(sources[0]) &&
		utils.HasDockerfile(sources[0]) {
		if err := copy.Copy(sources[0], workdir); err != nil {
			return err
		}
	} else {
		if err := packer.Pack(workdir, sources...); err != nil {
			return err
		}
		hooks := ctx.Get("hooks").([]*hook.Hook)
		for _, h := range hooks {
			if h.Name() == hook.HookNameBeforeBuild {
				if err := h.Run(workdir); err != nil {
					return err
				}
			}
		}
	}

	cloudType := ctx.Get("cloud_type").(string)
	name := ctx.Get("name").(string)
	if cloudType == types.CloudTypeK8S {
		data, err := packer.PackIntoK8SConfigMapFile(workdir)
		if err != nil {
			return err
		}
		ctx.Set("data", data)
	} else {
		docker := ctx.Get("docker").(containerruntimes.ContainerRuntime)
		if err := docker.BuildImage(ctx.GetContext(), workdir, name); err != nil {
			return err
		}
		nameWithTag := name + ":latest"
		if err := docker.TagImage(ctx.GetContext(), name, nameWithTag); err != nil {
			return err
		}

		ctx.Set("image", nameWithTag)
	}

	return nil
}
