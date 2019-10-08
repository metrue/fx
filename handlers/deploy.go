package handlers

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/apex/log"
	"github.com/metrue/fx/api"
	"github.com/metrue/fx/config"
	"github.com/metrue/fx/constants"
	"github.com/metrue/fx/deploy"
	dockerDeployer "github.com/metrue/fx/deploy/docker"
	k8sDeployer "github.com/metrue/fx/deploy/kubernetes"
	"github.com/metrue/fx/packer"
	"github.com/metrue/fx/utils"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

// Deploy deploy handle function
func Deploy(cfg config.Configer) HandleFunc {
	return func(ctx *cli.Context) (err error) {
		funcFile := ctx.Args().First()
		name := ctx.String("name")
		port := ctx.Int("port")
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

		workdir, err := ioutil.TempDir("/tmp", "fx-wd")
		if err != nil {
			return err
		}
		if err := packer.PackIntoDir(lang, string(body), workdir); err != nil {
			return err
		}

		var deployer deploy.Deployer
		if os.Getenv("KUBECONFIG") != "" {
			deployer, err = k8sDeployer.Create()
			if err != nil {
				return err
			}
		} else {
			bctx := context.Background()
			deployer, err = dockerDeployer.CreateClient(bctx)
			if err != nil {
				return err
			}
		}
		// TODO multiple ports support
		return deployer.Deploy(
			context.Background(),
			workdir,
			name,
			[]int32{int32(port)},
		)
	}
}
