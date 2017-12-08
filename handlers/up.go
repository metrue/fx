package handlers

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/gorilla/websocket"
	api "github.com/metrue/fx/docker-api"
	"github.com/metrue/fx/image"
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

func notify(connection *websocket.Conn, messageType int, message string) {
	err := connection.WriteMessage(messageType, []byte(message))
	if err != nil {
		log.Println("write: ", err)
	}
}

func cleanup(dir string) {
	format := "cleanup temp file [%s] error: %s\n"
	if err := os.RemoveAll(dir); err != nil {
		log.Printf(format, dir, err.Error())
	}
	dirTar := dir + ".tar"
	if err := os.RemoveAll(dirTar); err != nil {
		log.Printf(format, dirTar, err.Error())
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
	image.Get(dir, string(lang), body)

	api.Build(name, dir)
	notify(connection, messageType, "function built")
	api.Deploy(name, dir, strconv.Itoa(port))
	msg := fmt.Sprintf("function deployed at: %s:%s", utils.GetHostIP().String(), strconv.Itoa(port))
	notify(connection, messageType, msg)
	closeConnection(connection)
}
