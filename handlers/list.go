package handlers

import (
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/deploy"
	"github.com/metrue/fx/pkg/render"
	"github.com/metrue/fx/pkg/spinner"
)

// List command handle
func List(ctx context.Contexter) (err error) {
	const task = "deploying"
	spinner.Start(task)
	defer func() {
		spinner.Stop(task, err)
	}()

	cli := ctx.GetCliContext()
	deployer := ctx.Get("deployer").(deploy.Deployer)

	services, err := deployer.List(ctx.GetContext(), cli.Args().First())
	if err != nil {
		return err
	}

	render.Table(services)
	return nil
}
