package handlers

import (
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/deploy"
	"github.com/metrue/fx/pkg/spinner"
	"github.com/metrue/fx/utils"
)

// List command handle
func List() HandleFunc {
	return func(ctx *context.Context) (err error) {
		const task = "deploying"
		spinner.Start(task)
		defer func() {
			spinner.Stop(task, err)
		}()

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
