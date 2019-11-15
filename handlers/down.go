package handlers

import (
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/deploy"
	"github.com/metrue/fx/pkg/spinner"
)

// Down command handle
func Down() HandleFunc {
	return func(ctx *context.Context) (err error) {
		spinner.Start("deploying")
		defer spinner.Stop()

		cli := ctx.GetCliContext()
		services := cli.Args()
		runner := ctx.Get("deployer").(deploy.Deployer)
		for _, svc := range services {
			if err := runner.Destroy(ctx.Context, svc); err != nil {
				return err
			}
		}
		return nil
	}
}
