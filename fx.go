package main

import (
	"log"
	"os"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/google/uuid"
	engine "github.com/metrue/fx/api"
	"github.com/urfave/cli"
)

var api *engine.API

func init() {
	box := packr.NewBox("./api/images")
	api = engine.NewWithDockerRemoteAPI("127.0.0.1:1234", box)
}

func main() {
	app := cli.NewApp()
	app.Name = "fx"
	app.Usage = "make function as a service"
	app.Version = "0.2.2"

	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "initialize fx running enviroment",
			Action: func(c *cli.Context) error {
				return api.Init()
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
				return api.Up(name, c.Args().First())
			},
		},
		{
			Name:      "down",
			Usage:     "destroy a service",
			ArgsUsage: "[service 1, service 2, ....]",
			Action: func(c *cli.Context) error {
				return api.Down(c.Args())
			},
		},
		{
			Name:  "list",
			Usage: "list deployed services",
			Action: func(c *cli.Context) error {
				return api.List(c.Args().First())
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
				return api.Call(c.Args().First(), params)
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
