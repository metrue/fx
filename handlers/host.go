package handlers

import (
	"log"

	"github.com/metrue/fx/config"
	"github.com/metrue/fx/utils"
	"github.com/urfave/cli"
)

// AddHost add a host
func AddHost(cfg config.Configer) HandleFunc {
	return func(ctx *cli.Context) error {
		name := ctx.String("name")
		addr := ctx.String("host")
		user := ctx.String("user")
		password := ctx.String("password")
		host := config.NewHost(addr, user, password)
		if !host.Valid() {
			log.Fatalf("invaid host %v", host)
			return nil
		}

		if host.IsRemote() {
			if host.User == "" || host.Password == "" {
				log.Fatalf("the host to add is a remote, user and password for SSH login is required")
				return nil
			}
		}
		return cfg.AddMachine(name, host)
	}
}

// RemoveHost remove a host
func RemoveHost(cfg config.Configer) HandleFunc {
	return func(ctx *cli.Context) error {
		name := ctx.Args().First()
		if name == "" {
			log.Fatalf("no name given: fx infra remove <name>")
			return nil
		}
		return cfg.RemoveHost(name)
	}
}

// ListHosts list hosts
func ListHosts(cfg config.Configer) HandleFunc {
	return func(ctx *cli.Context) error {
		hosts, err := cfg.ListMachines()
		if err != nil {
			return err
		}

		return utils.OutputJSON(hosts)
	}
}
