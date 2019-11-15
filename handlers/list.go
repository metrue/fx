package handlers

import (
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/deploy"
	"github.com/metrue/fx/pkg/spinner"
	"github.com/metrue/fx/utils"
)

// List command handle
func List() HandleFunc {
	return func(ctx *context.Context) error {
		spinner.Start("deploying")
		defer spinner.Stop()

		cli := ctx.GetCliContext()
		deployer := ctx.Get("deployer").(deploy.Deployer)

		services, err := deployer.List(ctx.Context, cli.Args().First())
		if err != nil {
			return err
		}

		for _, service := range services {
			if err := utils.OutputJSON(service); err != nil {
				return err
			}
		}

		return nil
	}
}
