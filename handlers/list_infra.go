package handlers

import (
	"fmt"

	"github.com/metrue/fx/config"
	"github.com/metrue/fx/context"
)

// ListInfra list infra
func ListInfra(ctx *context.Context) (err error) {
	if _, err := config.Load(); err != nil {
		return err
	}

	conf, err := config.View()
	if err != nil {
		return err
	}
	fmt.Println(string(conf))
	return nil
}
