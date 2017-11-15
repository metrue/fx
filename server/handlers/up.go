package handlers

import (
	"log"
	"os"
	"path"
	"strconv"

	"github.com/gorilla/websocket"
	api "github.com/metrue/fx/server/docker-api"
	"github.com/metrue/fx/utils"
	"github.com/phayes/freeport"
	"github.com/rs/xid"
)

func closeConnection(connection *websocket.Conn) {
	connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "0"))
}

var funcNames = map[string]string{
	"go":     "/fx.go",
	"node":   "/fx.js",
	"ruby":   "/fx.rb",
	"python": "/fx.py",
}

func dispatchFuncion(lang string, data []byte, dir string) {
	fileName := funcNames[lang]
	f, err := os.Create(dir + fileName)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	n, err := f.Write(data)
	if err != nil {
		panic(err)
	}

	log.Println("func recved: %s", n)
}

func notify(connection *websocket.Conn, messageType int, message string) {
	err := connection.WriteMessage(messageType, []byte(message))
	if err != nil {
		log.Println("write: ", err)
	}
}

func initWorkDirectory(lang string, dir string) {
	err := utils.CopyDir(path.Join(os.Getenv("HOME"), ".fx/images", lang), dir)
	if err != nil {
		panic(err)
	}
}

func Up(
	lang []byte,
	body []byte,
	connection *websocket.Conn,
	messageType int,
) {
	// go func() {
	var guid = xid.New().String()
	var dir = guid
	var name = guid

	port, err := freeport.GetFreePort()
	if err != nil {
		panic(err)
	}

	initWorkDirectory(string(lang), dir)
	notify(connection, messageType, "work dir initialized")
	dispatchFuncion(string(lang), body, dir)
	notify(connection, messageType, "function dispatched")

	api.Build(name, dir)

	notify(connection, messageType, "function built")
	api.Deploy(name, dir, strconv.Itoa(port))
	notify(connection, messageType, "function deployed: http://localhost:"+strconv.Itoa(port))

	closeConnection(connection)
	// }()
}
