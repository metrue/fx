package handlers

import (
	"os"

	"github.com/apex/log"
	"github.com/metrue/fx/constants"
	dockerHTTP "github.com/metrue/fx/container_runtimes/docker/http"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/deploy"
	"github.com/metrue/fx/provision"
	"github.com/metrue/fx/types"
	"github.com/metrue/fx/utils"
)

// List command handle
func List() HandleFunc {
	return func(ctx *context.Context) (err error) {
		host := os.Getenv("DOCKER_REMOTE_HOST_ADDR")
		user := os.Getenv("DOCKER_REMOTE_HOST_USER")
		passord := os.Getenv("DOCKER_REMOTE_HOST_PASSWORD")

		cli := ctx.GetCliContext()
		services := []types.Service{}
		deployer := ctx.Get("deployer").(deploy.Deployer)

		// FIXME clean up following check, use deployer's List
		if host != "" && user != "" {
			provisioner := provision.NewWithHost(host, user, passord)
			if !provisioner.IsFxAgentRunning() {
				if err := provisioner.StartFxAgent(); err != nil {
					return err
				}
			}
			httpClient, err := dockerHTTP.Create(host, constants.AgentPort)
			if err != nil {
				return err
			}
			services, err = httpClient.ListContainer(cli.Args().First())
			if err != nil {
				log.Fatalf("list functions on machine %s failed: %v", host, err)
			}
		} else {
			services, err = deployer.List(ctx.Context, cli.Args().First())
			if err != nil {
				return err
			}
		}

		for _, service := range services {
			if err := utils.OutputJSON(service); err != nil {
				return err
			}
		}

		return nil
	}
}
