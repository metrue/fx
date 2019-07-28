package main

import (
	"fmt"
	"net"
	"os"
	"path"
	"strings"

	"github.com/apex/log"
	"github.com/gobuffalo/packr"
	"github.com/google/uuid"
	"github.com/metrue/fx/api"
	"github.com/metrue/fx/commands"
	"github.com/metrue/fx/config"
	"github.com/metrue/fx/constants"
	"github.com/metrue/fx/doctor"
	"github.com/metrue/fx/provision"
	"github.com/phayes/freeport"
	"github.com/urfave/cli"
)

func init() {
	if err := config.Init(path.Join(os.Getenv("HOME"), ".fx")); err != nil {
		log.Fatalf("Init config failed %s", err)
		os.Exit(1)
	}
}

func main() {
	app := cli.NewApp()
	app.Name = "fx"
	app.Usage = "makes function as a service"
	app.Version = "0.3.22"

	host, err := config.GetDefaultHost()
	if err != nil {
		fmt.Println(err)
		log.Fatalf("fx config is dirty, run 'fx autofix' to fix")
	}
	endpoint := net.JoinHostPort(host.Host, constants.AgentPort)
	box := packr.NewBox("./api/images")
	fx := api.NewWithDockerRemoteAPI(endpoint, box)

	app.Commands = []cli.Command{
		{
			Name:  "host",
			Usage: "manage hosts",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "add a new host",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "name, N",
							Usage: "a alias name for this host",
						},
						cli.StringFlag{
							Name:  "host, H",
							Usage: "host name or IP address of a host",
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
						name := c.String("name")
						host := c.String("host")
						user := c.String("user")
						password := c.String("password")
						return commands.AddHost(name, config.NewHost(host, user, password))
					},
				},
				{
					Name:  "remove",
					Usage: "remove an existing host",
					Action: func(c *cli.Context) error {
						if c.Args().First() == "" {
							log.Fatalf("no name given: fx host remove <host_name>")
							return nil
						}
						return commands.RemoveHost(c.Args().First())
					},
				},
				{
					Name:  "list",
					Usage: "list hosts",
					Action: func(c *cli.Context) error {
						return commands.ListHosts()
					},
				},
				{
					Name:  "default",
					Usage: "set/get default host",
					Action: func(c *cli.Context) error {
						if c.Args().First() != "" {
							return commands.SetDefaultHost(c.Args().First())
						}
						return commands.GetDefaultHost()
					},
				},
			},
		},
		{
			Name:  "doctor",
			Usage: "health check for fx",
			Action: func(c *cli.Context) error {
				host := c.String("host")
				user := c.String("user")
				password := c.String("password")
				return doctor.New(host, user, password).Start()
			},
		},
		{
			Name:  "provision",
			Usage: "initialize a host to be a fx server",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "host, H",
					Usage: "host name or IP address of a host",
					Value: "127.0.0.1",
				},
				cli.StringFlag{
					Name:  "user, U",
					Usage: "user name required for SSH login",
				},
				cli.StringFlag{
					Name:  "password, P",
					Usage: "password required for SSH login",
				},
				cli.StringFlag{
					Name:  "key, K",
					Usage: "full path to public key file",
				},
			},
			Action: func(c *cli.Context) error {
				host := c.String("host")
				user := c.String("user")
				password := c.String("password")
				opts := provision.Options{
					Host:     host,
					User:     user,
					Password: password,
				}
				provisionor := provision.New(opts)
				return provisionor.Start()
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
				cli.IntFlag{
					Name:  "port, p",
					Usage: "port number",
				},
			},
			Action: func(c *cli.Context) error {
				if err := fx.Health(); err != nil {
					log.Warn("fx is not healthy, make sure you have done 'fx provision' first")
					return nil
				}

				name := c.String("name")
				if name == "" {
					name = uuid.New().String()
				}
				port := c.Int("port")
				if port == 0 {
					freePort, err := freeport.GetFreePort()
					if err != nil {
						return err
					}
					port = freePort
				}
				return fx.Up(c.Args().First(), api.UpOptions{Name: name, Port: port})
			},
		},
		{
			Name:      "down",
			Usage:     "destroy a service",
			ArgsUsage: "[service 1, service 2, ....]",
			Action: func(c *cli.Context) error {
				if err := fx.Health(); err != nil {
					log.Warn("fx is not healthy, make sure you have done 'fx provision' first")
					return nil
				}

				return fx.Down(c.Args())
			},
		},
		{
			Name:  "list",
			Usage: "list deployed services",
			Action: func(c *cli.Context) error {
				if err := fx.Health(); err != nil {
					log.Warn("fx is not healthy, make sure you have done 'fx provision' first")
					return nil
				}

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
				if err := fx.Health(); err != nil {
					log.Warn("fx is not healthy, make sure you have done 'fx provision' first")
					return nil
				}

				params := strings.Join(c.Args()[1:], " ")
				return fx.Call(c.Args().First(), params)
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("fx startup with fatal: %v", err)
	}
}
