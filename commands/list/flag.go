package list

import (
	"fmt"
	"os"
	"flag"
	"net/url"
)

type argPtrs struct {
	addr *string
	help *bool
}

func setupFlags() (
	args *argPtrs,
	flagSet *flag.FlagSet,
){
	flagSet = flag.NewFlagSet("list", flag.ExitOnError)
	args = &argPtrs{
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

func parseArgs(
	s []string,
	ptrs *argPtrs,
	fs *flag.FlagSet,
) (funcs []string, addr string) {
	fs.Parse(s)
	if *(ptrs.help) {
		flagsAndExit(fs)
	}

	u := url.URL{
		Scheme: "ws",
		Host:   *(ptrs.addr),
		Path:   "/list",
	}
	addr = u.String()

	if fs.NFlag() == 0 {
		funcs = s
	} else {
		funcs = fs.Args()
	}
	return
}

func flagsAndExit(fs *flag.FlagSet) {
	fmt.Println("Flags:")
	fs.PrintDefaults()
	os.Exit(0)
}
