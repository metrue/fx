package handlers

import (
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/infra"
	"github.com/metrue/fx/pkg/render"
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
	render.Table([]types.Service{service})
	return nil
}
