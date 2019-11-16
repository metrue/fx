package handlers

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	containerruntimes "github.com/metrue/fx/container_runtimes"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/deploy"
	"github.com/metrue/fx/packer"
	"github.com/metrue/fx/pkg/spinner"
	"github.com/metrue/fx/types"
	"github.com/metrue/fx/utils"
	"github.com/pkg/errors"
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
func Up() HandleFunc {
	return func(ctx *context.Context) (err error) {
		spinner.Start("deploying")
		defer func() {
			spinner.Stop(err)
		}()

		cli := ctx.GetCliContext()
		funcFile := cli.Args().First()
		name := cli.String("name")
		port := cli.Int("port")

		if port < PortRange.min || port > PortRange.max {
			return fmt.Errorf("invalid port number: %d, port number should in range of %d -  %d", port, PortRange.min, PortRange.max)
		}

		lang := utils.GetLangFromFileName(funcFile)
		body, err := ioutil.ReadFile(funcFile)
		if err != nil {
			return errors.Wrap(err, "read source failed")
		}
		fn := types.Func{Language: lang, Source: string(body)}
		if os.Getenv("K3S") != "" {
			docker := ctx.Get("docker").(containerruntimes.ContainerRuntime)
			workdir := fmt.Sprintf("/tmp/fx-%d", time.Now().Unix())
			defer os.RemoveAll(workdir)

			if err := packer.PackIntoDir(fn, workdir); err != nil {
				return err
			}
			if err := docker.BuildImage(ctx.Context, workdir, name); err != nil {
				return err
			}
			if out, err := docker.PushImage(ctx.Context, name); err != nil {
				fmt.Println(out)
				return err
			}
		}

		deployer := ctx.Get("deployer").(deploy.Deployer)
		bindings := ctx.Get("bindings").([]types.PortBinding)
		return deployer.Deploy(
			ctx.Context,
			types.Func{Language: lang, Source: string(body)},
			name,
			bindings,
		)
	}
}
