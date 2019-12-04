package middlewares

import (
	"io/ioutil"

	"github.com/metrue/fx/context"
	"github.com/metrue/fx/types"
	"github.com/metrue/fx/utils"
	"github.com/pkg/errors"
)

// Parse parse input
func Parse(action string) func(ctx context.Contexter) (err error) {
	return func(ctx context.Contexter) error {
		cli := ctx.GetCliContext()
		switch action {
		case "up":
			funcFile := cli.Args().First()
			lang := utils.GetLangFromFileName(funcFile)
			body, err := ioutil.ReadFile(funcFile)
			if err != nil {
				return errors.Wrap(err, "read source failed")
			}
			fn := types.Func{
				Language: lang,
				Source:   string(body),
			}
			ctx.Set("fn", fn)

			name := cli.String("name")
			ctx.Set("name", name)
			port := cli.Int("port")
			ctx.Set("port", port)
		case "down":
			services := cli.Args()
			if len(services) == 0 {
				return errors.New("service name required")
			}
			svc := []string{}
			for _, service := range services {
				svc = append(svc, service)
			}
			ctx.Set("services", svc)
		case "list":
			name := cli.Args().First()
			ctx.Set("filter", name)
		}

		return nil
	}
}
