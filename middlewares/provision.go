package middlewares

import (
	"fmt"
	"strings"

	"github.com/metrue/fx/constants"
	dockerHTTP "github.com/metrue/fx/container_runtimes/docker/http"
	"github.com/metrue/fx/context"
	dockerInfra "github.com/metrue/fx/infra/docker"
	k8sInfra "github.com/metrue/fx/infra/k8s"
	"github.com/metrue/fx/types"
)

// Provision make sure infrastructure is healthy
func Provision(ctx context.Contexter) (err error) {
	host := ctx.Get("host").(string)
	port := ctx.Get("ssh_port").(string)
	keyfile := ctx.Get("ssh_key").(string)
	kubeconf := ctx.Get("kubeconf").(string)
	if host == "" && kubeconf == "" {
		return fmt.Errorf("at least host or kubeconf provided")
	}

	if host != "" {
		addr := strings.Split(host, "@")
		if len(addr) != 2 {
			return fmt.Errorf("invalid host information, should be format of <user>@<ip>")
		}
		user := addr[0]
		ip := addr[1]

		cloud, err := dockerInfra.Create(ip, user, port, keyfile)
		if err != nil {
			return err
		}
		ok, err := cloud.IsHealth()
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("target docker host is not healthy")
		}

		// TODO port should be configurable
		docker, err := dockerHTTP.Create(ip, constants.AgentPort)
		if err != nil {
			return err
		}
		ctx.Set("docker", docker)

		deployer, err := dockerInfra.CreateDeployer(docker)
		if err != nil {
			return err
		}
		ctx.Set("cloud_type", types.CloudTypeDocker)
		ctx.Set("deployer", deployer)
	} else if kubeconf != "" {
		deployer, err := k8sInfra.CreateDeployer(kubeconf)
		if err != nil {
			return err
		}
		ctx.Set("cloud_type", types.CloudTypeK8S)
		ctx.Set("deployer", deployer)
	}

	return nil
}
