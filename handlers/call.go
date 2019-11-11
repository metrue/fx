package handlers

import (
	"io/ioutil"
	"strings"

	"github.com/apex/log"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/packer"
	"github.com/metrue/fx/types"
	"github.com/metrue/fx/utils"
)

// Call command handle
func Call() HandleFunc {
	return func(ctx *context.Context) error {
		cli := ctx.GetCliContext()
		_ = strings.Join(cli.Args()[1:], " ")

		file := cli.Args().First()
		src, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalf("Read Source: %v", err)
			return err
		}
		log.Info("Read Source: \u2713")

		lang := utils.GetLangFromFileName(file)
		fn := types.Func{
			Language: lang,
			Source:   string(src),
		}
		if _, err := packer.Pack(file, fn); err != nil {
			panic(err)
		}

		// TODO not supported
		// if err := api.MustCreate(host.Host, constants.AgentPort).
		// 	Call(file, params, project); err != nil {
		// 	log.Fatalf("call functions on machine %s with %v failed: %v", name, params, err)
		// }

		return nil
	}
}
