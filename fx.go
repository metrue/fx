package main

import (
	"log"
	"os"

	"github.com/metrue/fx/commands"
	"github.com/metrue/fx/config"
	"github.com/metrue/fx/server"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "fx"
	app.Usage = "make function as a service"
	app.Version = "0.0.81"

	app.Commands = []cli.Command{
		{
			Name:  "serve",
			Usage: "start fx server on current host",
			Action: func(c *cli.Context) error {
				return server.Start(true)
			},
		},
		{
			Name:      "up",
			Usage:     "deploy a function or a group of functions",
			ArgsUsage: "[func.go func.js func.py func.rb ...]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "host, H",
					Usage: "fx server host, default is localhost",
				},
			},
			Action: func(c *cli.Context) error {
				host := c.String("host")
				if host == "" {
					host = config.GetGrpcEndpoint()
				}
				functionSources := c.Args()
				return commands.Up(host, functionSources)
			},
		},
		{
			Name:      "down",
			Usage:     "destroy a function or a group of functions",
			ArgsUsage: "[id1, id2, ...]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "host, H",
					Usage: "fx server host, default is localhost",
				},
			},
			Action: func(c *cli.Context) error {
				host := c.String("host")
				if host == "" {
					host = config.GetGrpcEndpoint()
				}
				functions := c.Args()
				return commands.Down(host, functions)
			},
		},
		{
			Name:  "list",
			Usage: "list deployed services",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "host, H",
					Usage: "fx server host, default is localhost",
				},
			},
			Action: func(c *cli.Context) error {
				host := c.String("host")
				if host == "" {
					host = config.GetGrpcEndpoint()
				}
				functions := c.Args()
				return commands.List(host, functions)
			},
		},
		{
			Name:  "use",
			Usage: "set target deploy server address, default is localhost",
			Action: func(c *cli.Context) error {
				return commands.Use(c.Args().First())
			},
		},
		{
			Name:  "status",
			Usage: "show fx status",
			Action: func(c *cli.Context) error {
				host := c.String("host")
				if host == "" {
					host = config.GetGrpcEndpoint()
				}
				return commands.Status(host)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
