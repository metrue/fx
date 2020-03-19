package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"regexp"

	"github.com/apex/log"
	"github.com/google/uuid"
	aurora "github.com/logrusorgru/aurora"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/handlers"
	"github.com/metrue/fx/middlewares"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

const version = "0.9.33"

func init() {
	go checkForUpdate()
}

func handle(fns ...func(ctx context.Contexter) error) func(ctx *cli.Context) error {
	return func(c *cli.Context) error {
		ctx := context.FromCliContext(c)
		for _, fn := range fns {
			if err := fn(ctx); err != nil {
				panic(err)
			}
		}
		return nil
	}
}

func checkForUpdate() {
	const releaseURL = "https://api.github.com/repos/metrue/fx/releases/latest"
	resp, err := http.Get(releaseURL)
	if err != nil {
		log.Debugf("Failed to fetch Github release page, error %v", err)
		return
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	var releaseJSON struct {
		Tag string `json:"tag_name"`
		URL string `json:"html_url"`
	}
	if err := decoder.Decode(&releaseJSON); err != nil {
		log.Debugf("Failed to decode Github release page JSON, error %v", err)
		return
	}
	if matched, err := regexp.MatchString(`^(\d+\.)(\d+\.)(\d+)$`, releaseJSON.Tag); err != nil || !matched {
		log.Debugf("Unofficial release %s?", releaseJSON.Tag)
		return
	}
	log.Debugf("Latest release tag is %s", releaseJSON.Tag)
	if releaseJSON.Tag != version {
		fmt.Fprintf(os.Stderr, "\nfx %s is available (you're using %s), get the latest release from: %s\n",
			releaseJSON.Tag, version, releaseJSON.URL)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "fx"
	app.Usage = "makes function as a service"
	app.Version = version

	defer func() {
		if r := recover(); r != nil {
			fmt.Println(aurora.Red("*****************"))
			fmt.Println(r)
			fmt.Println(aurora.Red("*****************"))
		}
	}()
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	defaultHost := user.Username + "@localhost"

	defaultSSHKeyFile, err := homedir.Expand("~/.ssh/id_rsa")
	if err != nil {
		panic(err)
	}

	app.Commands = []cli.Command{
		{
			Name:      "up",
			Usage:     "deploy a function",
			ArgsUsage: "[func.go func.js func.py func.rb ...]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Value: uuid.New().String(),
					Usage: "service name",
				},
				cli.IntFlag{
					Name:  "port, p",
					Usage: "port number",
				},
				cli.StringFlag{
					Name:  "host, H",
					Usage: "target host, <user>@<host>",
					Value: defaultHost,
				},
				cli.StringFlag{
					Name:  "ssh_port, P",
					Usage: "SSH port for target host",
					Value: "22",
				},
				cli.StringFlag{
					Name:  "ssh_key, K",
					Usage: "SSH key file for login target host",
					Value: defaultSSHKeyFile,
				},
				cli.StringFlag{
					Name:  "kubeconf, C",
					Usage: "kubeconf of kubernetes cluster",
				},
				cli.BoolFlag{
					Name:  "healthcheck, hc",
					Usage: "do a health check after service up",
				},
				cli.BoolFlag{
					Name:  "force, f",
					Usage: "force deploy a function or functions",
				},
			},
			Action: handle(
				middlewares.Parse("up"),
				middlewares.Language(),
				middlewares.Binding,
				middlewares.SSH,
				middlewares.Driver,
				middlewares.Build,
				handlers.Up,
			),
		},
		{
			Name:      "down",
			Usage:     "destroy a service",
			ArgsUsage: "[service 1, service 2, ....]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "ssh_port, P",
					Usage: "SSH port for target host",
					Value: "22",
				},
				cli.StringFlag{
					Name:  "ssh_key, K",
					Usage: "SSH key file for login target host",
					Value: defaultSSHKeyFile,
				},
				cli.StringFlag{
					Name:  "host, H",
					Usage: "target host, <user>@<host>",
					Value: defaultHost,
				},
				cli.StringFlag{
					Name:  "kubeconf, C",
					Usage: "kubeconf of kubernetes cluster",
				},
			},
			Action: handle(
				middlewares.Parse("down"),
				middlewares.SSH,
				middlewares.Driver,
				handlers.Down,
			),
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "list deployed services",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "format, f",
					Value: "table",
					Usage: "output format, 'table' and 'JSON' supported",
				},
				cli.StringFlag{
					Name:  "ssh_port, P",
					Usage: "SSH port for target host",
					Value: "22",
				},
				cli.StringFlag{
					Name:  "ssh_key, K",
					Usage: "SSH key file for login target host",
					Value: defaultSSHKeyFile,
				},
				cli.StringFlag{
					Name:  "host, H",
					Usage: "target host, <user>@<host>",
					Value: defaultHost,
				},
				cli.StringFlag{
					Name:  "kubeconf, C",
					Usage: "kubeconf of kubernetes cluster",
				},
			},
			Action: handle(
				middlewares.Parse("list"),
				middlewares.SSH,
				middlewares.Driver,
				handlers.List,
			),
		},
		{
			Name:  "image",
			Usage: "manage image of service",
			Subcommands: []cli.Command{
				{
					Name:  "build",
					Usage: "build a image",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "ssh_port, P",
							Usage: "SSH port for target host",
							Value: "22",
						},
						cli.StringFlag{
							Name:  "ssh_key, K",
							Usage: "SSH key file for login target host",
							Value: defaultSSHKeyFile,
						},
						cli.StringFlag{
							Name:  "host, H",
							Usage: "target host, <user>@<host>",
							Value: defaultHost,
						},
						cli.StringFlag{
							Name:  "kubeconf, C",
							Usage: "kubeconf of kubernetes cluster",
						},
						cli.StringFlag{
							Name:  "name, n",
							Usage: "image name",
							Value: uuid.New().String(),
						},
					},
					Action: handle(
						middlewares.Parse("image_build"),
						middlewares.Language(),
						middlewares.SSH,
						middlewares.Driver,
						middlewares.Build,
						handlers.BuildImage,
					),
				},
				{
					Name:  "export",
					Usage: "export the Docker project of service",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "output, o",
							Usage: "output directory",
						},
					},
					Action: handle(
						middlewares.Parse("image_export"),
						middlewares.Language(),
						handlers.ExportImage,
					),
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
