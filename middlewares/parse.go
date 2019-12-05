package middlewares

import (
	"os"

	"github.com/metrue/fx/context"
	"github.com/pkg/errors"
)

// Parse parse input
func Parse(action string) func(ctx context.Contexter) (err error) {
	return func(ctx context.Contexter) error {
		cli := ctx.GetCliContext()
		switch action {
		case "up":
			sources := []string{}
			for _, s := range cli.Args() {
				sources = append(sources, s)
			}
			if len(sources) == 0 {
				pwd, err := os.Getwd()
				if err != nil {
					return err
				}
				sources = append(sources, pwd)
			}
			ctx.Set("sources", sources)
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
