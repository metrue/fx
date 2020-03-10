package handlers

import (
	"fmt"
	"os"
	"time"

	"github.com/apex/log"
	"github.com/metrue/fx/bundle"
	"github.com/metrue/fx/constants"
	containerruntimes "github.com/metrue/fx/container_runtimes"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/hook"
	"github.com/metrue/fx/pkg/spinner"
)

// BuildImage build image
func BuildImage(ctx context.Contexter) (err error) {
	spinner.Start("building")
	defer func() {
		spinner.Stop("building", err)
	}()
	workdir := fmt.Sprintf("/tmp/fx-%d", time.Now().Unix())
	defer os.RemoveAll(workdir)

	fn := ctx.Get("fn").(string)
	deps := ctx.Get("deps").([]string)
	language := ctx.Get("language").(string)

	if err := bundle.Bundle(workdir, language, fn, deps...); err != nil {
		return err
	}
	if err := hook.RunBeforeBuildHook(workdir); err != nil {
		return err
	}

	docker := ctx.Get("docker").(containerruntimes.ContainerRuntime)
	nameWithTag := ctx.Get("tag").(string) + ":latest"
	if err := docker.BuildImage(ctx.GetContext(), workdir, nameWithTag); err != nil {
		return err
	}
	log.Infof("image built: %s %v", nameWithTag, constants.CheckedSymbol)
	return nil
}

// ExportImage export service's code into a directory
func ExportImage(ctx context.Contexter) (err error) {
	outputDir := ctx.Get("output").(string)
	fn := ctx.Get("fn").(string)
	deps := ctx.Get("deps").([]string)
	language := ctx.Get("language").(string)

	if err := bundle.Bundle(outputDir, language, fn, deps...); err != nil {
		return err
	}
	if err := hook.RunBeforeBuildHook(outputDir); err != nil {
		return err
	}

	log.Infof("exported to %v: %v", outputDir, constants.CheckedSymbol)
	return nil
}
