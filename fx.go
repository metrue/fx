package main

import (
	"fmt"
	"log"
	"os"

	"github.com/metrue/fx/server"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "fx"
	app.Usage = "make function as a service"
	app.Version = "0.0.2"

	app.Commands = []cli.Command{
		{
			Name:    "serve",
			Aliases: []string{"s"},
			Usage:   "manage fx server",
			Action: func(c *cli.Context) error {
				return server.Start(true)
			},
		},
		{
			Name:    "up",
			Aliases: []string{"u"},
			Usage:   "deploy a function or a group of functions",
			Action: func(c *cli.Context) error {
				fmt.Println("up: ", c.Args().First())
				return nil
			},
		},
		{
			Name:    "down",
			Aliases: []string{"d"},
			Usage:   "destroy a function or a group of functions",
			Action: func(c *cli.Context) error {
				fmt.Println("down: ", c.Args().First())
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list deployed services",
			Action: func(c *cli.Context) error {
				fmt.Println("list: ", c.Args().First())
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
