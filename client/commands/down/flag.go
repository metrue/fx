package down

import (
	"fmt"
	"os"
	"flag"
	"path"
	"io/ioutil"
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
	flagSet = flag.NewFlagSet("down", flag.ExitOnError)
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
	home := os.Getenv("HOME")
	configFile := path.Join(home, ".fx")
	buf, _ := ioutil.ReadFile(configFile)
	if len(buf) > 0 {
		fs.Set("addr", string(buf))
	}

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

func flagsAndExit(fs *flag.FlagSet) {
	fmt.Println("Flags:")
	fs.PrintDefaults()
	os.Exit(0)
}
