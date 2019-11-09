package middlewares

import (
	"os"

	"github.com/metrue/fx/constants"
	dockerHTTP "github.com/metrue/fx/container_runtimes/docker/http"
	dockerSDK "github.com/metrue/fx/container_runtimes/docker/sdk"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/deploy"
	dockerDeployer "github.com/metrue/fx/deploy/docker"
	k8sDeployer "github.com/metrue/fx/deploy/kubernetes"
)

// Setup create k8s or docker cli
func Setup(ctx *context.Context) (err error) {
	var deployer deploy.Deployer
	if os.Getenv("KUBECONFIG") != "" {
		deployer, err = k8sDeployer.Create()
		if err != nil {
			return err
		}
	} else {
		deployer, err = dockerDeployer.CreateClient(ctx.Context)
		if err != nil {
			return err
		}
	}
	ctx.Set("deployer", deployer)

	host := os.Getenv("DOCKER_REMOTE_HOST_ADDR")
	user := os.Getenv("DOCKER_REMOTE_HOST_USER")
	if host != "" && user != "" {
		httpClient, err := dockerHTTP.Create(host, constants.AgentPort)
		if err != nil {
			return err
		}
		ctx.Set("docker_http", httpClient)
	} else {
		cli, err := dockerSDK.CreateClient(ctx)
		if err != nil {
			return err
		}
		ctx.Set("docker_sdk", cli)
	}

	return nil
}
