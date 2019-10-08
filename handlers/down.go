package handlers

import (
	"github.com/apex/log"
	"github.com/metrue/fx/config"
	"github.com/metrue/fx/constants"
	api "github.com/metrue/fx/container_runtimes/docker/http"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

// Down command handle
func Down(cfg config.Configer) HandleFunc {
	return func(ctx *cli.Context) error {
		containerID := ctx.Args()
		hosts, err := cfg.ListActiveMachines()
		if err != nil {
			return errors.Wrapf(err, "list active machines failed: %v", err)
		}
		for name, host := range hosts {
			if err := api.MustCreate(host.Host, constants.AgentPort).
				Down(containerID); err != nil {
				return errors.Wrapf(err, "stop function on machine %s failed: %v", name, err)
			}
			log.Infof("stop function on machine %s: %v", name, constants.CheckedSymbol)
		}
		return nil
	}
}
