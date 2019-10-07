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
	"github.com/metrue/fx/container"
	dockerc "github.com/metrue/fx/container/docker"
	"github.com/metrue/fx/container/kubernetes"
	"github.com/metrue/fx/image/docker"
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

		// Build image
		wd, err := ioutil.TempDir("/tmp", "fx-wd")
		if err != nil {
			return err
		}
		if err := packer.PackIntoDir(lang, string(body), wd); err != nil {
			return err
		}
		imageBuilder, err := docker.CreateClient()
		if err != nil {
			return err
		}
		if err := imageBuilder.Build(wd, name); err != nil {
			return err
		}
		imageName := name

		var runner container.Runner
		if os.Getenv("KUBECONFIG") != "" {
			runner, err = kubernetes.Create()
			if err != nil {
				return err
			}

			imageName, err = imageBuilder.Push(name)
			if err != nil {
				return err
			}
		} else {
			runner, err = dockerc.CreateClient()
			if err != nil {
				return err
			}
		}

		// TODO multiple ports support
		if err := runner.Deploy(
			context.Background(),
			name,
			imageName,
			[]int32{int32(port)}); err != nil {
			return err
		}
		return nil
	}
}
