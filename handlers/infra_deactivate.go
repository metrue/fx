package handlers

import (
	"github.com/apex/log"
	"github.com/metrue/fx/config"
	"github.com/metrue/fx/constants"
	"github.com/urfave/cli"
)

// Deactivate a machine
func Deactivate(cfg config.Configer) HandleFunc {
	return func(ctx *cli.Context) error {
		name := ctx.Args().First()
		if name == "" {
			log.Fatalf("name required for: fx infra activate <name>")
			return nil
		}
		if err := cfg.DisableMachine(name); err != nil {
			log.Fatalf("could not disable %s: %v", name, err)
			return nil
		}
		log.Infof("machine %s deactive: %v", name, constants.CheckedSymbol)
		return nil
	}
}
