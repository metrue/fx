package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/apex/log"
	"github.com/google/uuid"
	aurora "github.com/logrusorgru/aurora"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/handlers"
	"github.com/metrue/fx/middlewares"
	"github.com/urfave/cli"
)

const version = "0.8.72"

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

	app.Commands = []cli.Command{
		{
			Name:  "infra",
			Usage: "manage infrastructure",
			Subcommands: []cli.Command{
				{
					Name:  "create",
					Usage: "create a infra for fx service",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "type, t",
							Usage: "infracture type, 'docker', 'k8s' and 'k3s' support",
						},
						cli.StringFlag{
							Name:  "name, n",
							Usage: "name to identify the infrastructure",
						},
						cli.StringFlag{
							Name:  "host",
							Usage: "user and ip of your host, eg. 'root@182.12.1.12'",
						},
						cli.StringFlag{
							Name:  "master",
							Usage: "serve as master node in K3S cluster, eg. 'root@182.12.1.12'",
						},
						cli.StringFlag{
							Name:  "agents",
							Usage: "serve as agent node in K3S cluster, eg. 'root@187.1. 2. 3,root@123.3.2.1'",
						},
					},

					Action: handle(
						middlewares.LoadConfig,
						handlers.Setup,
					),
				},
				{
					Name:  "list",
					Usage: "list all infrastructures",
					Action: handle(
						middlewares.LoadConfig,
						handlers.ListInfra,
					),
				},
				{
					Name:  "use",
					Usage: "set current context to target cloud with given name",
					Action: handle(
						middlewares.LoadConfig,
						handlers.UseInfra,
					),
				},
			},
		},
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
				middlewares.LoadConfig,
				middlewares.Provision,
				middlewares.Parse("up"),
				middlewares.Binding,
				middlewares.Build,
				handlers.Up,
			),
		},
		{
			Name:      "down",
			Usage:     "destroy a service",
			ArgsUsage: "[service 1, service 2, ....]",
			Action: handle(
				middlewares.Parse("down"),
				middlewares.LoadConfig,
				middlewares.Provision,
				handlers.Down,
			),
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "list deployed services",
			Action: handle(
				middlewares.Parse("list"),
				middlewares.LoadConfig,
				middlewares.Provision,
				handlers.List,
			),
		},
		{
			Name:  "call",
			Usage: "run a function instantly",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "host, H",
					Usage: "fx server host, default is localhost",
				},
			},
			Action: handle(handlers.Call),
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
							Name:  "tag, t",
							Usage: "image tag",
						},
					},
					Action: handle(
						middlewares.LoadConfig,
						middlewares.Provision,
						middlewares.Parse("image_build"),
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
						middlewares.LoadConfig,
						middlewares.Provision,
						middlewares.Parse("image_export"),
						handlers.ExportImage,
					),
				},
			},
		},
		{
			Name:   "doctor",
			Usage:  "health check for fx",
			Action: handle(handlers.Doctor),
		},
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
