package handlers

import (
	"fmt"
	"io/ioutil"

	"github.com/apex/log"
	"github.com/metrue/fx/config"
	"github.com/metrue/fx/constants"
	api "github.com/metrue/fx/container_runtimes/docker/http"
	"github.com/metrue/fx/packer"
	"github.com/metrue/fx/provision"
	"github.com/metrue/fx/types"
	"github.com/metrue/fx/utils"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

// PortRange usable port range https: //en.wikipedia.org/wiki/Ephemeral_port
var PortRange = struct {
	min int
	max int
}{
	min: 1023,
	max: 65535,
}

// Up command handle
func Up(cfg config.Configer) HandleFunc {
	return func(ctx *cli.Context) (err error) {
		funcFile := ctx.Args().First()
		name := ctx.String("name")
		port := ctx.Int("port")
		healtcheck := ctx.Bool("healthcheck")
		force := ctx.Bool("force")

		defer func() {
			if r := recover(); r != nil {
				log.Fatalf("fatal error happened: %v", r)
			}

			if err != nil {
				log.Fatalf("deploy function %s (%s) failed: %v", err)
			}
			log.Infof("function %s (%s) deployed successfully", name, funcFile)
		}()

		if port < PortRange.min || port > PortRange.max {
			return fmt.Errorf("invalid port number: %d, port number should in range of %d -  %d", port, PortRange.min, PortRange.max)
		}
		hosts, err := cfg.ListActiveMachines()
		if err != nil {
			return errors.Wrap(err, "list active machines failed")
		}

		if len(hosts) == 0 {
			log.Warnf("no active machines")
			return nil
		}

		// try to stop service firt
		if force {
			for n, host := range hosts {
				if err := api.MustCreate(host.Host, constants.AgentPort).
					Stop(name); err != nil {
					log.Infof("stop function %s on machine %s failed: %v", name, n, err)
				} else {
					log.Infof("stop function %s on machine %s: %v", name, n, constants.CheckedSymbol)
				}
			}
		}

		body, err := ioutil.ReadFile(funcFile)
		if err != nil {
			return errors.Wrap(err, "read source failed")
		}
		lang := utils.GetLangFromFileName(funcFile)

		fn := types.ServiceFunctionSource{
			Language: lang,
			Source:   string(body),
		}

		project, err := packer.Pack(name, fn)
		if err != nil {
			return errors.Wrapf(err, "could pack function %s (%s)", name, funcFile)
		}

		for n, host := range hosts {
			if !host.Provisioned {
				provisionor := provision.New(host)
				if err := provisionor.Start(); err != nil {
					return errors.Wrapf(err, "could not provision %s", n)
				}
				log.Infof("provision machine %v: %s", n, constants.CheckedSymbol)
				if err := cfg.UpdateProvisionedStatus(n, true); err != nil {
					return errors.Wrap(err, "update machine provision status failed")
				}
			}

			if err := api.MustCreate(host.Host, constants.AgentPort).
				Up(api.UpOptions{
					Body:       body,
					Lang:       lang,
					Name:       name,
					Port:       port,
					HealtCheck: healtcheck,
					Project:    project,
				}); err != nil {
				return errors.Wrapf(err, "up function %s(%s) to machine %s failed", name, funcFile, n)
			}
			log.Infof("up function %s(%s) to machine %s: %v", name, funcFile, n, constants.CheckedSymbol)
		}
		return nil
	}
}
