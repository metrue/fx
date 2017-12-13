package common

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/metrue/fx/config"
)

type argPtrs struct {
	addr *string
	help *bool
}

func SetupFlags(option string) (
	args *argPtrs,
	flagSet *flag.FlagSet,
) {

	uri := config.GrpcEndpoint
	if uri[:1] == ":" {
		uri = "localhost" + uri
	}

	flagSet = flag.NewFlagSet(option, flag.ExitOnError)
	args = &argPtrs{
		addr: flagSet.String(
			"addr",
			uri,
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

func ParseArgs(
	option string,
	s []string,
	ptrs *argPtrs,
	fs *flag.FlagSet,
) (funcs []string, addr string) {
	home := os.Getenv("HOME")
	configFile := path.Join(home, ".fx")
	buf, _ := ioutil.ReadFile(configFile)
	if len(buf) > 0 {
		fs.Set("addr", string(buf))
	}

	fs.Parse(s)

	if *(ptrs.help) {
		FlagsAndExit(fs)
	}

	addr = *(ptrs.addr)

	if fs.NFlag() == 0 {
		funcs = s
	} else {
		funcs = fs.Args()
	}
	return
}

func FlagsAndExit(fs *flag.FlagSet) {
	fmt.Println("Flags:")
	fs.PrintDefaults()
	os.Exit(0)
}
