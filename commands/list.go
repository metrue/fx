package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/metrue/fx/common"
)

// List lists all running function services
func List() {
	option := "list"
	nArgs := len(os.Args)
	args, flagSet := common.SetupFlags(option)
	if nArgs < 2 {
		common.FlagsAndExit(flagSet)
	}
	functions, address := common.ParseArgs(
		option,
		os.Args[2:],
		args,
		flagSet,
	)

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(address, nil)
	if common.CheckError(err) {
		return
	}
	defer conn.Close()

	err = conn.WriteMessage(
		websocket.TextMessage,
		[]byte(strings.Join(functions, " ")),
	)
	if common.CheckError(err) {
		return
	}

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
