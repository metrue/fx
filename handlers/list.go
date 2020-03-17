package handlers

import (
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/infra"
	"github.com/metrue/fx/pkg/renderrer"
)

// List command handle
func List(ctx context.Contexter) (err error) {
	cli := ctx.GetCliContext()
	format := ctx.Get("format").(string)

	for _, targetDeployer := range []string{"docker_deployer", "k8s_deployer"} {
		deployer, ok := ctx.Get(targetDeployer).(infra.Deployer)
		if !ok {
			continue
		}
		services, err := deployer.List(ctx.GetContext(), cli.Args().First())
		if err != nil {
			return err
		}
		if err := renderrer.Render(services, format); err != nil {
			return err
		}
	}
	return nil
}
