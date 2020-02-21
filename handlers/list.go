package handlers

import (
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/infra"
	"github.com/metrue/fx/pkg/renderrer"
)

// List command handle
func List(ctx context.Contexter) (err error) {
	cli := ctx.GetCliContext()
	deployer := ctx.Get("deployer").(infra.Deployer)
	format := ctx.Get("format").(string)

	services, err := deployer.List(ctx.GetContext(), cli.Args().First())
	if err != nil {
		return err
	}

	return renderrer.Render(services, format)
}
