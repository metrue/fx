package worker

import (
	"log"
	"os"
	"path"

	"../../utils"

	"github.com/gorilla/websocket"
)

var funcNames = map[string]string{
	"go":     "/fx.go",
	"node":   "/fx.js",
	"ruby":   "/fx.rb",
	"python": "/fx.py",
}

func initWorkDirectory(lang string, dir string) {
	err := utils.CopyDir(path.Join(os.Getenv("HOME"), ".fx/images", lang), dir)
	if err != nil {
		panic(err)
	}
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

func checkerror(err error) {
	if err != nil {
		panic(err)
	}
}

func notify(connection *websocket.Conn, messageType int, message string) {
	err := connection.WriteMessage(messageType, []byte(message))
	if err != nil {
		log.Println("write: ", err)
	}
}

func closeConnection(connection *websocket.Conn) {
	connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "0"))
}
