package handlers

import (
	"fmt"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/google/uuid"
	"github.com/metrue/fx/constants"
	containerruntimes "github.com/metrue/fx/container_runtimes"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/packer"
	"github.com/metrue/fx/utils"
	"github.com/otiai10/copy"
)

// BuildImage build image
func BuildImage(ctx context.Contexter) error {
	cli := ctx.GetCliContext()
	tag := cli.String("tag")
	if tag == "" {
		tag = uuid.New().String()
	}

	workdir := fmt.Sprintf("/tmp/fx-%d", time.Now().Unix())
	defer os.RemoveAll(workdir)

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
	}

	docker, ok := ctx.Get("docker").(containerruntimes.ContainerRuntime)
	if ok {
		nameWithTag := tag + ":latest"
		if err := docker.BuildImage(ctx.GetContext(), workdir, nameWithTag); err != nil {
			return err
		}
		log.Infof("image built: %v", constants.CheckedSymbol)
		return nil
	}
	return fmt.Errorf("no available docker cli")
}

// ExportImage export service's code into a directory
func ExportImage(ctx context.Contexter) (err error) {
	cli := ctx.GetCliContext()
	outputDir := cli.String("output")
	if outputDir == "" {
		log.Fatalf("output directory required")
		return nil
	}

	sources := ctx.Get("sources").([]string)

	if len(sources) == 0 {
		return fmt.Errorf("source file/directory of function required")
	}
	if len(sources) == 1 &&
		utils.IsDir(sources[0]) &&
		utils.HasDockerfile(sources[0]) {
		if err := copy.Copy(sources[0], outputDir); err != nil {
			return err
		}
	} else {
		if err := packer.Pack(outputDir, sources...); err != nil {
			return err
		}
	}

	log.Infof("exported to %v: %v", outputDir, constants.CheckedSymbol)
	return nil
}
