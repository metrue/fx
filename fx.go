package main

import (
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/metrue/fx/api"
	"github.com/urfave/cli"
)

var fx *api.API

func init() {
	version, err := api.Version(api.DockerRemoteAPIEndpoint)
	if err != nil {
		panic(err)
	}
	fx = api.NewWithDockerRemoteAPI(api.DockerRemoteAPIEndpoint, version)
}

func main() {
	app := cli.NewApp()
	app.Name = "fx"
	app.Usage = "makes function as a service"
	app.Version = "0.3.0"

	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "initialize fx running enviroment",
			Action: func(c *cli.Context) error {
				return fx.Init()
			},
		},
		{
			Name:      "up",
			Usage:     "deploy a function or a group of functions",
			ArgsUsage: "[func.go func.js func.py func.rb ...]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name, n",
					Usage: "service name",
				},
			},
			Action: func(c *cli.Context) error {
				name := c.String("name")
				if name == "" {
					name = uuid.New().String()
				}
				return fx.Up(name, c.Args().First())
			},
		},
		{
			Name:      "down",
			Usage:     "destroy a service",
			ArgsUsage: "[service 1, service 2, ....]",
			Action: func(c *cli.Context) error {
				return fx.Down(c.Args())
			},
		},
		{
			Name:  "list",
			Usage: "list deployed services",
			Action: func(c *cli.Context) error {
				return fx.List(c.Args().First())
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
				params := strings.Join(c.Args()[1:], " ")
				return fx.Call(c.Args().First(), params)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
