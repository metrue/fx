package handlers

import (
	"github.com/apex/log"
	"github.com/metrue/fx/config"
	"github.com/metrue/fx/constants"
	"github.com/metrue/fx/provision"
	"github.com/urfave/cli"
)

// Activate a machine to be fx server
func Activate(cfg config.Configer) HandleFunc {
	return func(ctx *cli.Context) error {
		name := ctx.Args().First()
		if name == "" {
			log.Fatalf("name required for: fx infra activate <name>")
			return nil
		}

		host, err := cfg.GetMachine(name)
		if err != nil {
			log.Fatalf("could get host %v, make sure you add it first", err)
			log.Info("You can add a machine by: \n fx infra add -Name <name> -H <ip or hostname> -U <user> -P <password>")
			return nil
		}
		if !host.Provisioned {
			provisionor := provision.New(host)
			if err := provisionor.Start(); err != nil {
				log.Fatalf("could not provision %s: %v", name, err)
				return nil
			}
			log.Infof("provision machine %v: %s", name, constants.CheckedSymbol)
			if err := cfg.UpdateProvisionedStatus(name, true); err != nil {
				log.Fatalf("update machine provision status failed: %v", err)
			}
		}

		if err := cfg.EnableMachine(name); err != nil {
			log.Fatalf("could not enable %s: %v", name, err)
			return nil
		}
		log.Infof("enble machine %v: %s", name, constants.CheckedSymbol)

		return nil
	}
}
