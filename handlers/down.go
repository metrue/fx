package handlers

import (
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/deploy"
	"github.com/metrue/fx/pkg/spinner"
)

// Down command handle
func Down(ctx *context.Context) (err error) {
	const task = "destroying"
	spinner.Start(task)
	defer func() {
		spinner.Stop(task, err)
	}()

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
