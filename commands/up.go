package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/metrue/fx/common"
	"github.com/metrue/fx/config"
)

type FunctionMeta struct {
	lang string
	path string
}

// Up starts the functions specified in flags
func Up() {
	option := "up"
	nArgs := len(os.Args)
	args, flagSet := common.SetupFlags(option)
	if nArgs == 2 {
		common.FlagsAndExit(flagSet)
	}
	functions, address := common.ParseArgs(
		option,
		os.Args[2:],
		args,
		flagSet,
	)

	fmt.Println("Deploy starting...")

	channel := make(chan bool)
	defer close(channel)

	numSuccess := 0
	numFail := 0

	for _, function := range functions {
		funcMeta := &FunctionMeta{
			lang: config.ExtLangMap[filepath.Ext(function)],
			path: function,
		}

		worker := NewWorker(funcMeta, address, channel)
		go worker.Work()
	}

	// Loop until all function deploy done
loop:
	for {
		select {
		case status := <-channel:
			if status {
				numSuccess++
			} else {
				numFail++
			}
		default:
			if numSuccess+numFail == len(functions) {
				fmt.Printf("Succed: %d\n", numSuccess)
				fmt.Printf("Failed: %d\n", numFail)
				fmt.Println("All deploy done!")
				break loop
			}
		}
	}
}
