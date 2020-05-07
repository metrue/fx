package middlewares

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/metrue/fx/constants"
	dockerHTTP "github.com/metrue/fx/container_runtimes/docker/http"
	"github.com/metrue/fx/context"
	dockerDriver "github.com/metrue/fx/driver/docker"
	k8sInfra "github.com/metrue/fx/driver/k8s"
	"github.com/metrue/fx/provisioner"
	"github.com/metrue/fx/provisioner/darwin"
	linux "github.com/metrue/fx/provisioner/linux"
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

			isRemote := (host != "127.0.0.1" && host != "localhost")
			hostOS := "linux"
			if isRemote {
				ok, err := sshClient.Connectable(provisioner.SSHConnectionTimeout)
				if err != nil {
					return err
				}
				if !ok {
					return fmt.Errorf("target host could not be connected with SSH")
				}

				var buf bytes.Buffer
				if err := sshClient.RunCommand("uname -a", ssh.CommandOptions{Stdout: &buf}); err != nil {
					return err
				}
				hostOS = buf.String()
			}
			if strings.Contains(hostOS, "darwin") {
				if err := darwin.New(sshClient).Provision(ctx.GetContext(), isRemote); err != nil {
					return err
				}
			} else {
				if err := linux.New(sshClient).Provision(ctx.GetContext(), isRemote); err != nil {
					return err
				}
			}
		}
		time.Sleep(2 * time.Second)
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
