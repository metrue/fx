package master

import (
	"fmt"
	"log"

	"../worker/uper"

	"github.com/gorilla/websocket"
)

func Up(functions []string, address string) {
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
		worker := uper.New(function, conn, channel)
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

func Down(functions []string, address string) {
	log.Print("Down functions: ", functions, " at address: ", address)
}
