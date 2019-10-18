package handlers

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/apex/log"
	"github.com/metrue/fx/config"
	"github.com/metrue/fx/constants"
	"github.com/metrue/fx/deploy"
	dockerDeployer "github.com/metrue/fx/deploy/docker"
	k8sDeployer "github.com/metrue/fx/deploy/kubernetes"
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

		body, err := ioutil.ReadFile(funcFile)
		if err != nil {
			return errors.Wrap(err, "read source failed")
		}
		lang := utils.GetLangFromFileName(funcFile)
		var deployer deploy.Deployer
		var bindings []types.PortBinding
		if os.Getenv("KUBECONFIG") != "" {
			deployer, err = k8sDeployer.Create()
			if err != nil {
				return err
			}
			bindings = []types.PortBinding{
				types.PortBinding{
					ServiceBindingPort:  80,
					ContainerExposePort: constants.FxContainerExposePort,
				},
				types.PortBinding{
					ServiceBindingPort:  443,
					ContainerExposePort: constants.FxContainerExposePort,
				},
			}
		} else {
			bctx := context.Background()
			deployer, err = dockerDeployer.CreateClient(bctx)
			if err != nil {
				return err
			}
			bindings = []types.PortBinding{
				types.PortBinding{
					ServiceBindingPort:  int32(port),
					ContainerExposePort: constants.FxContainerExposePort,
				},
			}
		}
		return deployer.Deploy(
			context.Background(),
			types.Func{Language: lang, Source: string(body)},
			name,
			bindings,
		)
	}
}
