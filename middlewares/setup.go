package middlewares

import (
	"fmt"
	"os"

	"github.com/metrue/fx/config"
	"github.com/metrue/fx/constants"
	dockerHTTP "github.com/metrue/fx/container_runtimes/docker/http"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/deploy"
	dockerDeployer "github.com/metrue/fx/deploy/docker"
	k3sDeployer "github.com/metrue/fx/deploy/k3s"
	k8sDeployer "github.com/metrue/fx/deploy/k8s"
	"github.com/metrue/fx/pkg/spinner"
)

// Setup create k8s or docker cli
func Setup(ctx *context.Context) (err error) {
	const task = "setup"
	spinner.Start(task)
	defer func() {
		spinner.Stop(task, err)
	}()

	fxConfig := ctx.Get("config").(*config.Config)
	cloud := fxConfig.Clouds[fxConfig.CurrentCloud]

	var deployer deploy.Deployer
	if cloud["type"] == config.CloudTypeDocker {
		docker, err := dockerHTTP.Create(cloud["host"], constants.AgentPort)
		if err != nil {
			return err
		}
		// TODO should clean up, but it needed in middlewares.Build
		ctx.Set("docker", docker)
		deployer, err = dockerDeployer.CreateClient(docker)
		if err != nil {
			return err
		}
	} else if cloud["type"] == config.CloudTypeK8S {
		if os.Getenv("K3S") != "" {
			deployer, err = k3sDeployer.Create()
			if err != nil {
				return err
			}
		} else if os.Getenv("KUBECONFIG") != "" {
			deployer, err = k8sDeployer.Create()
			if err != nil {
				return err
			}
		}
	} else {
		return fmt.Errorf("unsupport cloud type %s, please make sure you config is correct", cloud["type"])
	}

	ctx.Set("deployer", deployer)

	return nil
}
