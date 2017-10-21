package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
)

type upArgPtrs struct {
	addr *string
	help *bool
}

type downArgPtrs struct {
	addr *string
	help *bool
}

func setupUpFlags() (
	args *upArgPtrs,
	flagSet *flag.FlagSet,
) {
	flagSet = flag.NewFlagSet("up", flag.ExitOnError)
	args = &upArgPtrs{
		addr: flagSet.String(
			"addr",
			"localhost:8080",
			"Server address.",
		),
		help: flagSet.Bool(
			"help",
			false,
			"Help information.",
		),
	}
	return
}

func setupDownFlags() (
	args *downArgPtrs,
	flagSet *flag.FlagSet,
) {
	flagSet = flag.NewFlagSet("down", flag.ExitOnError)
	args = &downArgPtrs{
		addr: flagSet.String(
			"addr",
			"localhost:8080",
			"Server address.",
		),
		help: flagSet.Bool(
			"help",
			false,
			"Help information.",
		),
	}
	return
}

func parseUpArgs(
	s []string,
	ptrs *upArgPtrs,
	fs *flag.FlagSet,
) (funcs []string, addr string) {
	fs.Parse(s)
	if *(ptrs.help) {
		flagsAndExit(fs)
	}

	u := url.URL{
		Scheme: "ws",
		Host:   *(ptrs.addr),
		Path:   "/up",
	}
	addr = u.String()

	if fs.NFlag() == 0 {
		funcs = s
	} else {
		funcs = fs.Args()
	}
	return
}

func parseDownArgs(
	s []string,
	ptrs *downArgPtrs,
	fs *flag.FlagSet,
) (funcs []string, addr string) {
	fs.Parse(s)
	if *(ptrs.help) {
		flagsAndExit(fs)
	}

	u := url.URL{
		Scheme: "ws",
		Host:   *(ptrs.addr),
		Path:   "/down",
	}
	addr = u.String()

	if fs.NFlag() == 0 {
		funcs = s
	} else {
		funcs = fs.Args()
	}
	return
}

func checkMainFlag() {
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

func flagsAndExit(fs *flag.FlagSet) {
	fmt.Println("Flags:")
	fs.PrintDefaults()
	os.Exit(0)
}
