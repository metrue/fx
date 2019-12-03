package handlers

import (
	"fmt"

	"github.com/metrue/fx/config"
	"github.com/metrue/fx/context"
)

// ListInfra list infra
func ListInfra(ctx context.Contexter) (err error) {
	fxConfig := ctx.Get("config").(*config.Config)
	conf, err := fxConfig.View()
	if err != nil {
		return err
	}
	fmt.Println(string(conf))
	return nil
}
