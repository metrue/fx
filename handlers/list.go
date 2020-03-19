package handlers

import (
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/driver"
	"github.com/metrue/fx/pkg/renderrer"
)

// List command handle
func List(ctx context.Contexter) (err error) {
	cli := ctx.GetCliContext()
	format := ctx.Get("format").(string)

	for _, targetdriver := range []string{"docker_driver", "k8s_driver"} {
		driver, ok := ctx.Get(targetdriver).(driver.Driver)
		if !ok {
			continue
		}
		services, err := driver.List(ctx.GetContext(), cli.Args().First())
		if err != nil {
			return err
		}
		if err := renderrer.Render(services, format); err != nil {
			return err
		}
	}
	return nil
}
