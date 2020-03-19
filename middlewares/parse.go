package middlewares

import (
	"fmt"
	"os"
	"strings"

	"github.com/metrue/fx/context"
	"github.com/metrue/fx/utils"
	"github.com/urfave/cli"
)

type argsField struct {
	Type string
	Name string
	Env  string
}

func set(ctx context.Contexter, cli *cli.Context, fields []argsField) error {
	for _, f := range fields {
		if f.Type == "string" {
			if f.Name == "host" {
				addr := strings.Split(cli.String(f.Name), "@")
				fmt.Println("===>", cli.String(f.Name))
				if len(addr) != 2 {
					return fmt.Errorf("invalid host information, should be format of <user>@<ip>")
				}
				user := addr[0]
				ip := addr[1]
				ctx.Set("host", ip)
				ctx.Set("user", user)
			} else {

				ctx.Set(f.Name, cli.String(f.Name))
			}
		} else if f.Type == "int" {
			ctx.Set(f.Name, cli.Int(f.Name))
		} else if f.Type == "bool" {
			ctx.Set(f.Name, cli.Bool(f.Name))
		}

		if f.Env != "" && os.Getenv(f.Env) != "" {
			ctx.Set(f.Name, os.Getenv(f.Env))
		}
	}
	return nil
}

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

			if err := set(ctx, cli, []argsField{
				argsField{Name: "name", Type: "string"},
				argsField{Name: "port", Type: "int"},
				argsField{Name: "force", Type: "bool"},
				argsField{Name: "ssh_port", Type: "string", Env: "SSH_PORT"},
				argsField{Name: "ssh_key", Type: "string", Env: "SSH_KEY_FILE"},
				argsField{Name: "host", Type: "string", Env: "FX_HOST"},
				argsField{Name: "kubeconf", Type: "string", Env: "FX_KUBECONF"},
			}); err != nil {
				return err
			}

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

			if err := set(ctx, cli, []argsField{
				argsField{Name: "ssh_port", Type: "string", Env: "SSH_PORT"},
				argsField{Name: "ssh_key", Type: "string", Env: "SSH_KEY_FILE"},
				argsField{Name: "host", Type: "string", Env: "FX_HOST"},
				argsField{Name: "kubeconf", Type: "string", Env: "FX_KUBECONF"},
			}); err != nil {
				return err
			}

		case "list":
			name := cli.Args().First()
			ctx.Set("filter", name)
			format := cli.String("format")
			ctx.Set("format", format)
			if err := set(ctx, cli, []argsField{
				argsField{Name: "ssh_port", Type: "string", Env: "SSH_PORT"},
				argsField{Name: "ssh_key", Type: "string", Env: "SSH_KEY_FILE"},
				argsField{Name: "host", Type: "string", Env: "FX_HOST"},
				argsField{Name: "kubeconf", Type: "string", Env: "FX_KUBECONF"},
			}); err != nil {
				return err
			}

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
			if err := set(ctx, cli, []argsField{
				argsField{Name: "tag", Type: "string"},
				argsField{Name: "ssh_port", Type: "string", Env: "SSH_PORT"},
				argsField{Name: "ssh_key", Type: "string", Env: "SSH_KEY_FILE"},
				argsField{Name: "host", Type: "string", Env: "FX_HOST"},
				argsField{Name: "kubeconf", Type: "string", Env: "FX_KUBECONF"},
			}); err != nil {
				return err
			}

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
