package commands

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gorilla/websocket"
	"github.com/metrue/fx/common"
	"github.com/metrue/fx/utils"
)

type FunctionMeta struct {
	Lang    string
	Path    string
	Content string
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

	var funcList []FunctionMeta
	for _, function := range functions {
		data, err := ioutil.ReadFile(function)
		if err != nil {
			panic(err)
		}

		funcMeta := FunctionMeta{
			Lang:    utils.GetLangFromFileName(function),
			Path:    function,
			Content: string(data),
		}
		funcList = append(funcList, funcMeta)
	}
	funcListData, jsonErr := json.Marshal(funcList)
	if jsonErr != nil {
		panic(jsonErr)
	}

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(address, nil)
	if err != nil {
		panic(err)
	}

	err = conn.WriteMessage(websocket.TextMessage, funcListData)
	if err != nil {
		panic(err)
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		panic(err)
	}
	fmt.Println(string(msg))
}
