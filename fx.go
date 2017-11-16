package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/metrue/fx/commands/down"
	"github.com/metrue/fx/commands/list"
	"github.com/metrue/fx/commands/up"
	"github.com/metrue/fx/server"
)

const version string = "0.0.2"

const usage = `Usage:
  $ fx up   func1 func2 ...       deploy a function or a group of functions
  $ fx down func1 func2 ...       destroy a function or a group of functions
  $ fx list                       list deployed services
  $ fx server                     manage fx server
  $ fx --version                  show current version of f(x)
`

func versionAndExit() {
	fmt.Println(version)
	os.Exit(0)
}

func helpAndExit() {
	fmt.Print(usage)
	os.Exit(0)
}

func checkFlag() {
	helpPtr := flag.Bool(
		"help",
		false,
		"Help information.",
	)
	versionPtr := flag.Bool(
		"version",
		false,
		"Version information.",
	)

	flag.Parse()
	if *helpPtr {
		helpAndExit()
	}
	if *versionPtr {
		versionAndExit()
	}
}

func main() {
	nArgs := len(os.Args)
	if nArgs < 2 {
		helpAndExit()
	}
	checkFlag()

	switch os.Args[1] {
	case "server":
		server.Start()
	case "up":
		up.Up()
	case "down":
		down.Down()
	case "list":
		list.List()
	default:
		fmt.Print(usage)
		os.Exit(1)
	}
}
