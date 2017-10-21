package main

import (
	"fmt"
	"os"
	// "os/signal"

	"./master"
)

const version string = "0.0.2"

const usage = `Usage:
  $ fx up   func1 func2 ...       deploy a function or a group of functions
  $ fx down func1 func2 ...       destroy a function or a group of functions
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

func main() {
	nArgs := len(os.Args)
	if nArgs < 2 {
		helpAndExit()
	}
	checkMainFlag()

	upArgs, upFlagSet := setupUpFlags()
	downArgs, downFlagSet := setupDownFlags()
	listArgs, listFlagSet := setupListFlags()

	switch os.Args[1] {
	case "up":
		if nArgs == 2 {
			flagsAndExit(upFlagSet)
		}
		functions, address := parseUpArgs(
			os.Args[2:],
			upArgs,
			upFlagSet,
		)
		master.Up(functions, address)
	case "down":
		if nArgs == 2 {
			flagsAndExit(downFlagSet)
		}
		functions, address := parseDownArgs(
			os.Args[2:],
			downArgs,
			downFlagSet,
		)
		master.Up(functions, address)
	case "list":
		functions, address := parseListArgs(
			os.Args[2:],
			listArgs,
			listFlagSet,
		)
		master.List(functions, address)
	default:
		fmt.Print(usage)
		os.Exit(1)
	}
}
