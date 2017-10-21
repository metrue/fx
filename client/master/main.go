package master

import (
	"fmt"
	"log"
	"strings"

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
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(address, nil)
	if checkErrPrint(err) { return }
	defer conn.Close()

	err = conn.WriteMessage(
		websocket.TextMessage,
		[]byte(strings.Join(functions, " ")),
	)
	if checkErrPrint(err) { return }

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, 1000) {
				return
			}
			fmt.Println(err)
			if websocket.IsUnexpectedCloseError(err, 1000) {
				return
			}
			continue
		}
		fmt.Println(string(msg))
	}
}

func List(functions []string, address string) {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(address, nil)
	if checkErrPrint(err) { return }
	defer conn.Close()

	err = conn.WriteMessage(
		websocket.TextMessage,
		[]byte(strings.Join(functions, " ")),
	)
	if checkErrPrint(err) { return }

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, 1000) {
				return
			}
			fmt.Println(err)
			if websocket.IsUnexpectedCloseError(err, 1000) {
				return
			}
			continue
		}
		fmt.Println(string(msg))
	}
}

func checkErrPrint(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}
