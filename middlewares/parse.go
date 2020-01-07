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
			sources := []string{}
			for _, s := range cli.Args() {
				if utils.IsDir(s) || utils.IsRegularFile(s) {
					sources = append(sources, s)
				} else {
					return fmt.Errorf("no such file or directory: %s", s)
				}
			}
			ctx.Set("sources", sources)
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
		case "image_build":
			sources := []string{}
			for _, s := range cli.Args() {
				sources = append(sources, s)
			}
			ctx.Set("sources", sources)
			tag := cli.String("tag")
			if tag == "" {
				tag = uuid.New().String()
			}
			ctx.Set("tag", tag)
		case "image_export":
			sources := []string{}
			for _, s := range cli.Args() {
				sources = append(sources, s)
			}
			ctx.Set("sources", sources)
			outputDir := cli.String("output")
			if outputDir == "" {
				return fmt.Errorf("output directory required")
			}
			ctx.Set("output", outputDir)
		}

		return nil
	}
}
