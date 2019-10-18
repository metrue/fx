package handlers

import (
	"context"
	"os"

	"github.com/metrue/fx/config"
	"github.com/metrue/fx/deploy"
	dockerDeployer "github.com/metrue/fx/deploy/docker"
	k8sDeployer "github.com/metrue/fx/deploy/kubernetes"
	"github.com/urfave/cli"
)

// Down command handle
func Down(cfg config.Configer) HandleFunc {
	return func(ctx *cli.Context) (err error) {
		services := ctx.Args()
		c := context.Background()
		var runner deploy.Deployer
		if os.Getenv("KUBECONFIG") != "" {
			runner, err = k8sDeployer.Create()
			if err != nil {
				return err
			}
		} else {
			runner, err = dockerDeployer.CreateClient(c)
			if err != nil {
				return err
			}
		}
		for _, svc := range services {
			if err := runner.Destroy(c, svc); err != nil {
				return err
			}
		}
		return nil
	}
}
