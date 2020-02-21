package handlers

import (
	"github.com/apex/log"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/infra"
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
	deployer := ctx.Get("deployer").(infra.Deployer)
	bindings := ctx.Get("bindings").([]types.PortBinding)
	force := ctx.Get("force").(bool)
	if force && name != "" {
		if err := deployer.Destroy(ctx.GetContext(), name); err != nil {
			log.Warnf("destroy service %s failed: %v", name, err)
		}
	}

	if err := deployer.Deploy(
		ctx.GetContext(),
		fn,
		name,
		image,
		bindings,
	); err != nil {
		return err
	}

	service, err := deployer.GetStatus(ctx.GetContext(), name)
	if err != nil {
		return err
	}
	return renderrer.Render([]types.Service{service}, "table")
}
