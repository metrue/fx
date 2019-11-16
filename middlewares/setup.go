package middlewares

import (
	"os"

	"github.com/metrue/fx/constants"
	containerruntimes "github.com/metrue/fx/container_runtimes"
	dockerHTTP "github.com/metrue/fx/container_runtimes/docker/http"
	dockerSDK "github.com/metrue/fx/container_runtimes/docker/sdk"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/deploy"
	dockerDeployer "github.com/metrue/fx/deploy/docker"
	k3sDeployer "github.com/metrue/fx/deploy/k3s"
	k8sDeployer "github.com/metrue/fx/deploy/kubernetes"
	"github.com/metrue/fx/pkg/spinner"
	"github.com/metrue/fx/provision"
)

// Setup create k8s or docker cli
func Setup(ctx *context.Context) (err error) {
	const task = "setup"
	spinner.Start(task)
	defer func() {
		spinner.Stop(task, err)
	}()

	host := os.Getenv("DOCKER_REMOTE_HOST_ADDR")
	user := os.Getenv("DOCKER_REMOTE_HOST_USER")
	passord := os.Getenv("DOCKER_REMOTE_HOST_PASSWORD")
	var docker containerruntimes.ContainerRuntime
	if host != "" && user != "" {
		provisioner := provision.NewWithHost(host, user, passord)
		if !provisioner.IsFxAgentRunning() {
			if err := provisioner.StartFxAgent(); err != nil {
				return err
			}
		}

		docker, err = dockerHTTP.Create(host, constants.AgentPort)
		if err != nil {
			return err
		}
	} else {
		docker, err = dockerSDK.CreateClient(ctx)
		if err != nil {
			return err
		}
	}
	ctx.Set("docker", docker)

	var deployer deploy.Deployer
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
	} else {
		deployer, err = dockerDeployer.CreateClient(ctx.Context)
		if err != nil {
			return err
		}
	}
	ctx.Set("deployer", deployer)

	return nil
}
