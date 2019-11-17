package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"

	"github.com/apex/log"
	"github.com/google/uuid"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/handlers"
	"github.com/metrue/fx/middlewares"
	"github.com/urfave/cli"
)

const version = "0.8.3"

func init() {
	go checkForUpdate()
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

	app.Commands = []cli.Command{
		{
			Name:  "init",
			Usage: "start fx agent on host",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "master",
					Usage: "master node",
				},
				cli.StringFlag{
					Name:  "agents",
					Usage: "agent nodes",
				},
				cli.StringFlag{
					Name:  "user",
					Usage: "user acount name for SSH login",
				},
				cli.StringFlag{
					Name:  "password",
					Usage: "password for SSH login",
				},
			},
			Action: func(c *cli.Context) error {
				return handlers.Init()(context.FromCliContext(c))
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
			Action: func(c *cli.Context) error {
				ctx := context.FromCliContext(c)
				if err := ctx.Use(middlewares.Setup); err != nil {
					log.Fatalf("%v", err)
				}
				if err := ctx.Use(middlewares.Binding); err != nil {
					log.Fatalf("%v", err)
				}
				if err := ctx.Use(middlewares.Parse); err != nil {
					log.Fatalf("%v", err)
				}
				if err := ctx.Use(middlewares.Build); err != nil {
					log.Fatalf("%v", err)
				}
				return handlers.Up()(ctx)
			},
		},
		{
			Name:      "down",
			Usage:     "destroy a service",
			ArgsUsage: "[service 1, service 2, ....]",
			Action: func(c *cli.Context) error {
				ctx := context.FromCliContext(c)
				if err := ctx.Use(middlewares.Setup); err != nil {
					log.Fatalf("%v", err)
				}
				return handlers.Down()(ctx)
			},
		},
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "list deployed services",
			Action: func(c *cli.Context) error {
				ctx := context.FromCliContext(c)
				if err := ctx.Use(middlewares.Setup); err != nil {
					log.Fatalf("%v", err)
				}
				return handlers.List()(ctx)
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
				return handlers.Call()(context.FromCliContext(c))
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
						ctx := context.FromCliContext(c)
						if err := ctx.Use(middlewares.Setup); err != nil {
							log.Fatalf("%v", err)
						}
						return handlers.BuildImage()(ctx)
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
						return handlers.ExportImage()(context.FromCliContext(c))
					},
				},
			},
		},
		{
			Name:  "doctor",
			Usage: "health check for fx",
			Action: func(c *cli.Context) error {
				return handlers.Doctor()(context.FromCliContext(c))
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
