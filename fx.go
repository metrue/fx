package main

import (
	"os"
	"path"

	"github.com/apex/log"
	"github.com/google/uuid"
	"github.com/metrue/fx/config"
	"github.com/metrue/fx/handlers"
	"github.com/urfave/cli"
)

var cfg *config.Config

func init() {
	configDir := path.Join(os.Getenv("HOME"), ".fx")
	cfg := config.New(configDir)

	if err := cfg.Init(); err != nil {
		log.Fatalf("Init config failed %s", err)
		os.Exit(1)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "fx"
	app.Usage = "makes function as a service"
	app.Version = "0.6.2"

	app.Commands = []cli.Command{
		{
			Name:  "infra",
			Usage: "manage infrastructure of fx",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "add a new machine",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, N",
							Usage: "a alias name for this machine",
						},
						cli.StringFlag{
							Name:  "host, H",
							Usage: "host name or IP address of a machine",
						},
						cli.StringFlag{
							Name:  "user, U",
							Usage: "user name required for SSH login",
						},
						cli.StringFlag{
							Name:  "password, P",
							Usage: "password required for SSH login",
						},
					},
					Action: func(c *cli.Context) error {
						return handlers.AddHost(cfg)(c)
					},
				},
				{
					Name:  "remove",
					Usage: "remove an existing machine",
					Action: func(c *cli.Context) error {
						return handlers.RemoveHost(cfg)(c)
					},
				},
				{
					Name:    "list",
					Aliases: []string{"ls"},
					Usage:   "list machines",
					Action: func(c *cli.Context) error {
						return handlers.ListHosts(cfg)(c)
					},
				},
				{
					Name:  "activate",
					Usage: "enable a machine be a host of fx infrastructure",
					Action: func(c *cli.Context) error {
						return handlers.Activate(cfg)(c)
					},
				},
				{
					Name:  "deactivate",
					Usage: "disable a machine be a host of fx infrastructure",
					Action: func(c *cli.Context) error {
						return handlers.Deactivate(cfg)(c)
					},
				},
			},
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
					Action: func(c *cli.Context) error {
						return handlers.BuildImage(cfg)(c)
					},
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
					Action: func(c *cli.Context) error {
						return handlers.ExportImage()(c)
					},
				},
			},
		},
		{
			Name:  "doctor",
			Usage: "health check for fx",
			Action: func(c *cli.Context) error {
				return handlers.Doctor(cfg)(c)
			},
		},
		{
			Name:      "up",
			Usage:     "deploy a function or a group of functions",
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
			Action: func(c *cli.Context) error {
				return handlers.Up(cfg)(c)
			},
		},
		{
			Name:      "deploy",
			Usage:     "deploy a function or a group of functions",
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
			Action: func(c *cli.Context) error {
				return handlers.Deploy(cfg)(c)
			},
		},
		{
			Name:      "down",
			Usage:     "destroy a service",
			ArgsUsage: "[service 1, service 2, ....]",
			Action: func(c *cli.Context) error {
				return handlers.Down(cfg)(c)
			},
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "list deployed services",
			Action: func(c *cli.Context) error {
				return handlers.List(cfg)(c)
			},
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
			Action: func(c *cli.Context) error {
				return handlers.Call(cfg)(c)
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("fx startup with fatal: %v", err)
	}
}
