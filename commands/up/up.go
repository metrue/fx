package up

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gorilla/websocket"
	"github.com/metrue/fx/commands/common"
	"github.com/metrue/fx/config"
)

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
	dialer := websocket.Dialer{}

	channel := make(chan bool)
	defer close(channel)

	numSuccess := 0
	numFail := 0

	for _, function := range functions {
		conn, _, err := dialer.Dial(address, nil)
		if err != nil {
			log.Print(err)
			numFail++
			continue
		}
		lang := config.ExtLangMap[filepath.Ext(function)]
		worker := NewWorker(function, lang, conn, channel)
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
