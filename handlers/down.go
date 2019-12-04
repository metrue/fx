package handlers

import (
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/deploy"
)

// Down command handle
func Down(ctx context.Contexter) (err error) {
	services := ctx.Get("services").([]string)
	runner := ctx.Get("deployer").(deploy.Deployer)
	for _, svc := range services {
		if err := runner.Destroy(ctx.GetContext(), svc); err != nil {
			return err
		}
	}
	return nil
}
