package middlewares

import (
	"fmt"
	"os"

	"github.com/metrue/fx/config"
	"github.com/metrue/fx/constants"
	dockerHTTP "github.com/metrue/fx/container_runtimes/docker/http"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/infra"
	dockerInfra "github.com/metrue/fx/infra/docker"
	k8sInfra "github.com/metrue/fx/infra/k8s"
	"github.com/pkg/errors"
)

// Setup create k8s or docker cli
func Setup(ctx context.Contexter) (err error) {
	fxConfig := ctx.Get("config").(*config.Config)
	cloud := fxConfig.Clouds[fxConfig.CurrentCloud]

	var deployer infra.Deployer
	if cloud["type"] == config.CloudTypeDocker {
		docker, err := dockerHTTP.Create(cloud["host"], constants.AgentPort)
		if err != nil {
			return errors.Wrapf(err, "please make sure docker is installed and running on your host")
		}
		// TODO should clean up, but it needed in middlewares.Build
		ctx.Set("docker", docker)
		deployer, err = dockerInfra.CreateDeployer(docker)
		if err != nil {
			return err
		}
	} else if cloud["type"] == config.CloudTypeK8S {
		if os.Getenv("KUBECONFIG") != "" {
			deployer, err = k8sInfra.CreateDeployer(cloud["kubeconfig"])
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
