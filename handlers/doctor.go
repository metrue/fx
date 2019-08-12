package handlers

import (
	"github.com/apex/log"
	"github.com/metrue/fx/config"
	"github.com/metrue/fx/constants"
	"github.com/metrue/fx/doctor"
	"github.com/urfave/cli"
)

// Doctor command handle
func Doctor(cfg config.Configer) HandleFunc {
	return func(ctx *cli.Context) error {
		hosts, err := cfg.ListMachines()
		if err != nil {
			log.Fatalf("list machines failed %v", err)
			return nil
		}
		for name, h := range hosts {
			if err := doctor.New(h).Start(); err != nil {
				log.Warnf("machine %s is in dirty state: %v", name, err)
			} else {
				log.Infof("machine %s is in healthy state: %s", name, constants.CheckedSymbol)
			}
		}
		return nil
	}
}
