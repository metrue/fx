package middlewares

import (
	"fmt"
	"strings"

	"github.com/metrue/fx/constants"
	dockerHTTP "github.com/metrue/fx/container_runtimes/docker/http"
	"github.com/metrue/fx/context"
	dockerInfra "github.com/metrue/fx/infra/docker"
	k8sInfra "github.com/metrue/fx/infra/k8s"
)

// Driver initialize infrastructure driver
func Driver(ctx context.Contexter) (err error) {
	host := ctx.Get("host").(string)
	kubeconf := ctx.Get("kubeconf").(string)
	if host != "" {
		addr := strings.Split(host, "@")
		if len(addr) != 2 {
			return fmt.Errorf("invalid host information, should be format of <user>@<ip>")
		}
		ip := addr[1]
		// TODO port should be configurable
		docker, err := dockerHTTP.Create(ip, constants.AgentPort)
		if err != nil {
			return err
		}

		driver, err := dockerInfra.CreateDeployer(docker)
		if err != nil {
			return err
		}
		ctx.Set("docker_driver", driver)
	}

	if kubeconf != "" {
		driver, err := k8sInfra.CreateDeployer(kubeconf)
		if err != nil {
			return err
		}
		ctx.Set("k8s_driver", driver)
	}

	return nil
}
