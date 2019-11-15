package handlers

import (
	"fmt"
	"io/ioutil"

	"github.com/apex/log"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/deploy"
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
		defer spinner.Stop()

		cli := ctx.GetCliContext()
		funcFile := cli.Args().First()
		name := cli.String("name")
		port := cli.Int("port")

		defer func() {
			if r := recover(); r != nil {
				log.Fatalf("fatal error happened: %v", r)
			}

			if err != nil {
				log.Fatalf("deploy function %s (%s) failed: %v", err)
			}
			log.Infof("function %s (%s) deployed successfully", name, funcFile)
		}()

		if port < PortRange.min || port > PortRange.max {
			return fmt.Errorf("invalid port number: %d, port number should in range of %d -  %d", port, PortRange.min, PortRange.max)
		}

		body, err := ioutil.ReadFile(funcFile)
		if err != nil {
			return errors.Wrap(err, "read source failed")
		}
		lang := utils.GetLangFromFileName(funcFile)
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
