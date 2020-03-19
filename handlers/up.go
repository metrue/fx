package handlers

import (
	"github.com/apex/log"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/driver"
	"github.com/metrue/fx/pkg/renderrer"
	"github.com/metrue/fx/types"
)

// Up command handle
func Up(ctx context.Contexter) (err error) {
	fn, ok := ctx.Get("data").(string)
	if !ok {
		fn = ""
	}
	image, ok := ctx.Get("image").(string)
	if !ok {
		image = ""
	}
	name := ctx.Get("name").(string)
	bindings := ctx.Get("bindings").([]types.PortBinding)
	force := ctx.Get("force").(bool)

	for _, targetdriver := range []string{"docker_driver", "k8s_driver"} {
		driver, ok := ctx.Get(targetdriver).(driver.Driver)
		if !ok {
			continue
		}
		if force && name != "" {
			if err := driver.Destroy(ctx.GetContext(), name); err != nil {
				log.Warnf("destroy service %s failed: %v", name, err)
			}
		}

		if err := driver.Deploy(
			ctx.GetContext(),
			fn,
			name,
			image,
			bindings,
		); err != nil {
			return err
		}

		service, err := driver.GetStatus(ctx.GetContext(), name)
		if err != nil {
			return err
		}
		if err := renderrer.Render([]types.Service{service}, "table"); err != nil {
			return err
		}
	}
	return nil
}
