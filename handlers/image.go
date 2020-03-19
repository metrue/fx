package handlers

import (
	"github.com/apex/log"
	"github.com/metrue/fx/bundle"
	"github.com/metrue/fx/constants"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/hook"
)

// BuildImage build image
func BuildImage(ctx context.Contexter) (err error) {
	image := ctx.Get("image").(string)
	log.Infof("image built: %s %v", image, constants.CheckedSymbol)
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
