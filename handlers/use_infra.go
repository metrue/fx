package handlers

import (
	"github.com/metrue/fx/config"
	"github.com/metrue/fx/context"
)

// UseInfra use infra
func UseInfra(ctx *context.Context) error {
	fxConfig := ctx.Get("config").(*config.Config)
	cli := ctx.GetCliContext()
	return fxConfig.Use(cli.Args().First())
}
