package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/metrue/fx/commands"
	"github.com/metrue/fx/server"
)

const version string = "0.0.2"

const usage = `Usage:
  $ fx up   func1 func2 ...       deploy a function or a group of functions
  $ fx down func1 func2 ...       destroy a function or a group of functions
  $ fx list                       list deployed services
  $ fx serve                      manage fx server
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

func parseFlags() (version bool, help bool, verbose bool) {
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

	verbosePtr := flag.Bool(
		"verbose",
		false,
		"Make the operation more talkative",
	)

	flag.Parse()

	version = *versionPtr
	help = *helpPtr
	verbose = *verbosePtr

	return version, help, verbose
}

func main() {
	nArgs := len(os.Args)
	if nArgs < 2 {
		helpAndExit()
	}

	version, help, verbose := parseFlags()
	if help {
		helpAndExit()
	}

	if version {
		versionAndExit()
	}

	switch os.Args[1] {
	case "serve":
		server.Start(verbose)
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
