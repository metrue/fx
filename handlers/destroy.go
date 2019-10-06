package handlers

import (
	"context"
	"os"

	"github.com/metrue/fx/config"
	"github.com/metrue/fx/container"
	dockerc "github.com/metrue/fx/container/docker"
	"github.com/metrue/fx/container/kubernetes"
	"github.com/urfave/cli"
)

// Destroy command handle
func Destroy(cfg config.Configer) HandleFunc {
	return func(ctx *cli.Context) (err error) {
		services := ctx.Args()
		var runner container.Runner
		if os.Getenv("KUBECONFIG") != "" {
			runner, err = kubernetes.Create()
			if err != nil {
				return err
			}
		} else {
			runner, err = dockerc.CreateClient()
			if err != nil {
				return err
			}
		}
		for _, svc := range services {
			if err := runner.Destroy(context.Background(), svc); err != nil {
				return err
			}
		}
		return nil
	}
}
