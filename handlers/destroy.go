package handlers

import (
	"context"
	"fmt"
	"os"

	"github.com/metrue/fx/config"
	"github.com/metrue/fx/container/kubernetes"
	"github.com/urfave/cli"
)

// Destroy command handle
func Destroy(cfg config.Configer) HandleFunc {
	return func(ctx *cli.Context) error {
		services := ctx.Args()
		if os.Getenv("KUBECONFIG") != "" {
			runner, err := kubernetes.Create()
			if err != nil {
				return err
			}
			for _, svc := range services {
				if err := runner.Destroy(context.Background(), svc); err != nil {
					return err
				}
			}
		} else {
			return fmt.Errorf("no KUBECONFIG set in environment variables")
		}
		return nil
	}
}
