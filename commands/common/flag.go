package common

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"
)

type argPtrs struct {
	addr *string
	help *bool
}

func SetupFlags(option string) (
	args *argPtrs,
	flagSet *flag.FlagSet,
) {
	flagSet = flag.NewFlagSet(option, flag.ExitOnError)
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

	u := url.URL{
		Scheme: "ws",
		Host:   *(ptrs.addr),
		Path:   "/" + option,
	}
	addr = u.String()

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
