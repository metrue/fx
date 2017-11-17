package down

import (
	"fmt"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/metrue/fx/commands/common"
)

func Down() {
	option := "down"
	nArgs := len(os.Args)
	args, flagSet := common.SetupFlags(option)
	if nArgs == 2 {
		common.FlagsAndExit(flagSet)
	}
	functions, address := common.parseArgs(
		option,
		os.Args[2:],
		args,
		flagSet,
	)

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(address, nil)
	if checkErr(err) {
		return
	}
	defer conn.Close()

	err = conn.WriteMessage(
		websocket.TextMessage,
		[]byte(strings.Join(functions, " ")),
	)
	if checkErr(err) {
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

func checkErr(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}
