package handlers

import (
	"io/ioutil"
	"strings"

	"github.com/apex/log"
	"github.com/metrue/fx/config"
	"github.com/metrue/fx/constants"
	api "github.com/metrue/fx/container_runtimes/docker/http"
	"github.com/metrue/fx/packer"
	"github.com/metrue/fx/types"
	"github.com/metrue/fx/utils"
	"github.com/urfave/cli"
)

// Call command handle
func Call(cfg config.Configer) HandleFunc {
	return func(ctx *cli.Context) error {
		params := strings.Join(ctx.Args()[1:], " ")
		hosts, err := cfg.ListActiveMachines()
		if err != nil {
			log.Fatalf("list active machines failed: %v", err)
		}

		file := ctx.Args().First()
		src, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalf("Read Source: %v", err)
			return err
		}
		log.Info("Read Source: \u2713")

		lang := utils.GetLangFromFileName(file)
		fn := types.ServiceFunctionSource{
			Language: lang,
			Source:   string(src),
		}
		project, err := packer.Pack(file, fn)
		if err != nil {
			panic(err)
		}

		for name, host := range hosts {
			if err := api.MustCreate(host.Host, constants.AgentPort).
				Call(file, params, project); err != nil {
				log.Fatalf("call functions on machine %s with %v failed: %v", name, params, err)
			}
		}
		return nil
	}
}
