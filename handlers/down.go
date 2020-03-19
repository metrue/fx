package handlers

import (
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/driver"
)

// Down command handle
func Down(ctx context.Contexter) (err error) {
	services := ctx.Get("services").([]string)
	for _, targetdriver := range []string{"docker_driver", "k8s_driver"} {
		driver, ok := ctx.Get(targetdriver).(driver.Driver)
		if !ok {
			continue
		}
		for _, svc := range services {
			if err := driver.Destroy(ctx.GetContext(), svc); err != nil {
				return err
			}
		}
	}

	return nil
}
