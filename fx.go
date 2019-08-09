package main

import (
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

func fx(host config.Host) *api.API {
	box := packr.NewBox("./api/images")
	fx := api.New(cfg, box)
	if err := fx.Init(host); err != nil {
		log.Fatalf("Could not finish fx initialization: %v", err)
	}
	return fx
}

func main() {
	app := cli.NewApp()
	app.Name = "fx"
	app.Usage = "makes function as a service"
	app.Version = "0.4.0"

	commander := commands.New(cfg)

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
						name := c.String("name")
						host := c.String("host")
						user := c.String("user")
						password := c.String("password")
						return commander.AddHost(name, config.NewHost(host, user, password))
					},
				},
				{
					Name:  "remove",
					Usage: "remove an existing machine",
					Action: func(c *cli.Context) error {
						if c.Args().First() == "" {
							log.Fatalf("no name given: fx infra remove <name>")
							return nil
						}
						return commander.RemoveHost(c.Args().First())
					},
				},
				{
					Name:  "list",
					Usage: "list machines",
					Action: func(c *cli.Context) error {
						return commander.ListHosts()
					},
				},
				{
					Name:  "activate",
					Usage: "enable a machine be a host of fx infrastructure",
					Action: func(c *cli.Context) error {
						name := c.Args().First()
						if name == "" {
							log.Fatalf("name required for: fx infra activate <name>")
							return nil
						}

						host, err := cfg.GetMachine(name)
						if err != nil {
							log.Fatalf("could get host %v, make sure you add it first", err)
							log.Info("You can add a machine by: \n fx infra add -Name <name> -H <ip or hostname> -U <user> -P <password>")
							return nil
						}
						provisionor := provision.New(host)
						if err := provisionor.Start(); err != nil {
							log.Fatalf("could not provision %s: %v", name, err)
							return nil
						}
						log.Infof("provision machine %v: %s", name, constants.CheckedSymbol)

						if err := cfg.EnableMachine(name); err != nil {
							log.Fatalf("could not enable %s: %v", name, err)
							return nil
						}
						log.Infof("enble machine %v: %s", name, constants.CheckedSymbol)

						return nil
					},
				},
				{
					Name:  "deactivate",
					Usage: "disable a machine be a host of fx infrastructure",
					Action: func(c *cli.Context) error {
						// TODO set it status to be inactive
						return nil
					},
				},
			},
		},
		{
			Name:  "doctor",
			Usage: "health check for fx",
			Action: func(c *cli.Context) error {
				hosts, err := cfg.ListMachines()
				if err != nil {
					log.Fatalf("list machines failed %v", err)
					return nil
				}
				for name, h := range hosts {
					if err := doctor.New(h).Start(); err != nil {
						log.Warnf("machine %s is in dirty state: %v", name, err)
					} else {
						log.Infof("machine %s is in healthy state: %s", name, constants.CheckedSymbol)
					}
				}
				return nil
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
				name := c.String("name")
				if name == "" {
					name = uuid.New().String()
				}
				port := c.Int("port")
				if port == 0 {
					log.Fatalf("invalid port %d", port)
					return nil
				}
				hosts, err := cfg.ListActiveMachines()
				if err != nil {
					log.Fatalf("list active machines failed: %v", err)
				}
				for n, host := range hosts {
					src := c.Args().First()
					if err := fx(host).Up(src, api.UpOptions{Name: name, Port: port}); err != nil {
						log.Fatalf("up function %s(%s) to machine %s failed: %v", name, src, n, err)
					} else {
						log.Fatalf("up function %s(%s) to machine %s: %v", name, src, n, constants.CheckedSymbol)
					}
				}
				return nil
			},
		},
		{
			Name:      "down",
			Usage:     "destroy a service",
			ArgsUsage: "[service 1, service 2, ....]",
			Action: func(c *cli.Context) error {
				hosts, err := cfg.ListActiveMachines()
				if err != nil {
					log.Fatalf("list active machines failed: %v", err)
				}
				for name, host := range hosts {
					if err := fx(host).Down(c.Args()); err != nil {
						log.Fatalf("stop function on machine %s failed: %v", name, err)
					} else {
						log.Fatalf("stop function on machine %s: %v", name, constants.CheckedSymbol)
					}
				}
				return nil
			},
		},
		{
			Name:  "list",
			Usage: "list deployed services",
			Action: func(c *cli.Context) error {
				hosts, err := cfg.ListActiveMachines()
				if err != nil {
					log.Fatalf("list active machines failed: %v", err)
				}
				for name, host := range hosts {
					if err := fx(host).List(c.Args().First()); err != nil {
						log.Fatalf("list functions on machine %s failed: %v", name, err)
					} else {
						log.Fatalf("list functions on machine %s: %v", name, constants.CheckedSymbol)
					}
				}
				return nil
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
				hosts, err := cfg.ListActiveMachines()
				if err != nil {
					log.Fatalf("list active machines failed: %v", err)
				}
				for name, host := range hosts {
					if err := fx(host).Call(c.Args().First(), params); err != nil {
						log.Fatalf("call functions on machine %s with %v failed: %v", name, params, err)
					} else {
						log.Infof("call functions on machine %s with %v: %v", name, params, constants.CheckedSymbol)
					}
				}
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatalf("fx startup with fatal: %v", err)
	}
}
