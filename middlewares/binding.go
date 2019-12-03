package middlewares

import (
	"fmt"
	"os"

	"github.com/metrue/fx/constants"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/types"
	"github.com/phayes/freeport"
)

// PortRange usable port range https: //en.wikipedia.org/wiki/Ephemeral_port
var PortRange = struct {
	min int
	max int
}{
	min: 1023,
	max: 65535,
}

// Binding create bindings
func Binding(ctx context.Contexter) (err error) {
	port := ctx.Get("port").(int)
	if port == 0 {
		port, err = freeport.GetFreePort()
		if err != nil {
			return err
		}
	}
	if port < PortRange.min || port > PortRange.max {
		return fmt.Errorf("invalid port number: %d, port number should in range of %d -  %d", port, PortRange.min, PortRange.max)
	}

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
