package handlers

import (
	"github.com/apex/log"
	"github.com/metrue/fx/config"
	"github.com/metrue/fx/constants"
	api "github.com/metrue/fx/container_runtimes/docker/http"
	"github.com/urfave/cli"
)

// List command handle
func List(cfg config.Configer) HandleFunc {
	return func(ctx *cli.Context) error {
		hosts, err := cfg.ListActiveMachines()
		if err != nil {
			log.Fatalf("list active machines failed: %v", err)
		}
		for name, host := range hosts {
			if err := api.MustCreate(host.Host, constants.AgentPort).List(ctx.Args().First()); err != nil {
				log.Fatalf("list functions on machine %s failed: %v", name, err)
			}
		}
		return nil
	}
}
