package handlers

import (
	"fmt"

	"github.com/metrue/fx/context"
	"github.com/metrue/fx/deploy"
	"github.com/metrue/fx/pkg/spinner"
	"github.com/metrue/fx/types"
)

// PortRange usable port range https: //en.wikipedia.org/wiki/Ephemeral_port
var PortRange = struct {
	min int
	max int
}{
	min: 1023,
	max: 65535,
}

// Up command handle
func Up(ctx *context.Context) (err error) {
	const task = "deploying"
	spinner.Start(task)
	defer func() {
		spinner.Stop(task, err)
	}()

	cli := ctx.GetCliContext()
	name := cli.String("name")
	port := cli.Int("port")

	if port < PortRange.min || port > PortRange.max {
		return fmt.Errorf("invalid port number: %d, port number should in range of %d -  %d", port, PortRange.min, PortRange.max)
	}

	fn := ctx.Get("fn").(types.Func)
	image := ctx.Get("image").(string)
	deployer := ctx.Get("deployer").(deploy.Deployer)
	bindings := ctx.Get("bindings").([]types.PortBinding)
	return deployer.Deploy(
		ctx.Context,
		fn,
		name,
		image,
		bindings,
	)
}
