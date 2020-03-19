package middlewares

import (
	"fmt"
	"time"

	"github.com/apex/log"
	"github.com/metrue/fx/constants"
	dockerHTTP "github.com/metrue/fx/container_runtimes/docker/http"
	"github.com/metrue/fx/context"
	dockerDriver "github.com/metrue/fx/driver/docker"
	k8sInfra "github.com/metrue/fx/driver/k8s"
	"github.com/metrue/fx/provisioners"
	"github.com/metrue/go-ssh-client"
)

// Driver initialize infrastructure driver
func Driver(ctx context.Contexter) (err error) {
	host := ctx.Get("host").(string)
	sshClient := ctx.Get("ssh").(ssh.Clienter)
	kubeconf := ctx.Get("kubeconf").(string)
	if host != "" {
		// TODO port should be configurable
		docker := dockerHTTP.New(host, constants.AgentPort)
		driver := dockerDriver.New(dockerDriver.Options{
			DockerClient: docker,
		})

		if err := driver.Ping(ctx.GetContext()); err != nil {
			log.Infof("provisioning %s ...", host)

			provisioner := provisioners.New(sshClient)
			isRemote := (host != "127.0.0.1" && host != "localhost")
			if err := provisioner.Provision(ctx.GetContext(), isRemote); err != nil {
				return err
			}
			time.Sleep(2 * time.Second)
		}
		if err := docker.Initialize(); err != nil {
			return fmt.Errorf("initialize docker client failed: %s", err)
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
