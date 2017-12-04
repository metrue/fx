package handlers

import (
	"fmt"
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
	connection.WriteMessage(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, "0"),
	)
}

var funcNames = map[string]string{
	"go":     "/fx.go",
	"node":   "/fx.js",
	"ruby":   "/fx.rb",
	"python": "/fx.py",
	"php":    "/fx.php",
	"julia":  "/fx.jl",
	"java":   "/src/main/java/fx/Fx.java",
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

	log.Println("func recved:", n)
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

func cleanup(dir string) {
	if err := os.RemoveAll(dir); err != nil {
		log.Printf("cleanup [%s] error: %s\n", dir, err.Error())
	}
	dirTar := dir + ".tar"
	if err := os.RemoveAll(dirTar); err != nil {
		log.Printf("cleanup [%s] error: %s\n", dirTar, err.Error())
	}
}

// Up spins up a new function
func Up(lang []byte, body []byte, connection *websocket.Conn, messageType int) {
	port, err := freeport.GetFreePort()
	if err != nil {
		panic(err)
	}
	var guid = xid.New().String()
	var dir = path.Join(os.TempDir(), "fx-", guid)
	defer cleanup(dir)
	var name = guid
	initWorkDirectory(string(lang), dir)
	notify(connection, messageType, "work dir initialized")
	dispatchFuncion(string(lang), body, dir)
	notify(connection, messageType, "function dispatched")
	api.Build(name, dir)
	notify(connection, messageType, "function built")
	api.Deploy(name, dir, strconv.Itoa(port))
	msg := fmt.Sprintf("function deployed at: %s:%s", utils.GetHostIP().String(), strconv.Itoa(port))
	notify(connection, messageType, msg)
	closeConnection(connection)
}
