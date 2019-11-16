package middlewares

import (
	"os"

	"github.com/metrue/fx/constants"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/types"
)

// Binding create bindings
func Binding(ctx *context.Context) error {
	cli := ctx.GetCliContext()
	port := cli.Int("port")

	var bindings []types.PortBinding
	if os.Getenv("KUBECONFIG") != "" {
		bindings = []types.PortBinding{
			types.PortBinding{
				ServiceBindingPort:  80,
				ContainerExposePort: constants.FxContainerExposePort,
			},
			types.PortBinding{
				ServiceBindingPort:  443,
				ContainerExposePort: constants.FxContainerExposePort,
			},
			types.PortBinding{
				ServiceBindingPort:  int32(port),
				ContainerExposePort: constants.FxContainerExposePort,
			},
		}
	} else {
		bindings = []types.PortBinding{
			types.PortBinding{
				ServiceBindingPort:  int32(port),
				ContainerExposePort: constants.FxContainerExposePort,
			},
		}
	}
	ctx.Set("bindings", bindings)
	return nil
}
