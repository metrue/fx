package middlewares

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/utils"
)

// Parse parse input
func Parse(action string) func(ctx context.Contexter) (err error) {
	return func(ctx context.Contexter) error {
		cli := ctx.GetCliContext()
		switch action {
		case "up":
			if !cli.Args().Present() {
				return fmt.Errorf("no function given")
			}

			if !utils.IsRegularFile(cli.Args().First()) {
				return fmt.Errorf("invalid function source file: %s", cli.Args().First())
			}
			ctx.Set("fn", cli.Args().First())

			deps := []string{}
			for ind, s := range cli.Args() {
				if ind != 0 {
					deps = append(deps, s)
				}
			}
			ctx.Set("deps", deps)

			name := cli.String("name")
			ctx.Set("name", name)
			port := cli.Int("port")
			ctx.Set("port", port)
			force := cli.Bool("force")
			ctx.Set("force", force)
		case "down":
			services := cli.Args()
			if len(services) == 0 {
				return fmt.Errorf("service name required")
			}
			svc := []string{}
			for _, service := range services {
				svc = append(svc, service)
			}
			ctx.Set("services", svc)
		case "list":
			name := cli.Args().First()
			ctx.Set("filter", name)
			format := cli.String("format")
			ctx.Set("format", format)
		case "image_build":
			if !cli.Args().Present() {
				return fmt.Errorf("no function given")
			}

			if !utils.IsRegularFile(cli.Args().First()) {
				return fmt.Errorf("invalid function source file: %s", cli.Args().First())
			}
			ctx.Set("fn", cli.Args().First())

			deps := []string{}
			for ind, s := range cli.Args() {
				if ind != 0 {
					deps = append(deps, s)
				}
			}
			ctx.Set("deps", deps)

			tag := cli.String("tag")
			if tag == "" {
				tag = uuid.New().String()
			}
			ctx.Set("tag", tag)
		case "image_export":
			if !cli.Args().Present() {
				return fmt.Errorf("no function given")
			}

			if !utils.IsRegularFile(cli.Args().First()) {
				return fmt.Errorf("invalid function source file: %s", cli.Args().First())
			}
			ctx.Set("fn", cli.Args().First())

			deps := []string{}
			for ind, s := range cli.Args() {
				if ind != 0 {
					deps = append(deps, s)
				}
			}
			ctx.Set("deps", deps)

			outputDir := cli.String("output")
			if outputDir == "" {
				return fmt.Errorf("output directory required")
			}
			ctx.Set("output", outputDir)
		}

		return nil
	}
}
