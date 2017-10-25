package up

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gorilla/websocket"
)

var langs = map[string]string{
	".js": "node",
	".go": "go",
	".rb": "ruby",
	".py": "python",
}

func Up() {
	fmt.Println("Deploy starting...")

	nArgs := len(os.Args)
	args, flagSet := setupFlags()
	if nArgs == 2 {
		flagsAndExit(flagSet)
	}
	functions, address := parseArgs(
		os.Args[2:],
		args,
		flagSet,
	)

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
		lang := langs[filepath.Ext(function)]
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
				if numSuccess + numFail == len(functions) {
					fmt.Printf("Succed: %d\n", numSuccess)
					fmt.Printf("Failed: %d\n", numFail)
					fmt.Println("All deploy done!")
					break loop
				}
			}
		}
}
