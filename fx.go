package main

import (
	"flag"
	"fmt"
	"fx/commands"
	"os"

	"github.com/metrue/fx/server"
)

const version string = "0.0.2"

const usage = `Usage:
  $ fx up   func1 func2 ...       deploy a function or a group of functions
  $ fx down func1 func2 ...       destroy a function or a group of functions
  $ fx list                       list deployed services
  $ fx serve                     manage fx server
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
	case "serve":
		server.Start()
	case "up":
		commands.Up()
	case "down":
		commands.Down()
	case "list":
		commands.List()
	default:
		fmt.Print(usage)
		os.Exit(1)
	}
}
