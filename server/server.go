package server

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/metrue/fx/common"
	"github.com/metrue/fx/config"
	"github.com/metrue/fx/env"
	"github.com/metrue/fx/handlers"
	Message "github.com/metrue/fx/message"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "I am OK, %s!", r.URL.Path[1:])
}

func up(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade: ", err)
	}

	defer c.Close()

	mt, data, err := c.ReadMessage()
	if err != nil {
		log.Printf("read error: %s", err.Error())
		return
	}
	var funcList []common.FunctionMeta
	json.Unmarshal(data, &funcList)

	count := len(funcList)
	upResultCh := make(chan Message.UpMsgMeta, count)
	for _, funcMeta := range funcList {
		go handlers.Up(funcMeta, upResultCh)
	}

	// collect down result
	var ups []Message.UpMsgMeta
	for upResult := range upResultCh {
		ups = append(ups, upResult)
		if len(ups) == count {
			close(upResultCh)
		}
	}

	msg := Message.CreateUpMessage(ups)
	c.WriteMessage(mt, []byte(msg))
	closeConn(c, "0")
}

func list(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error during upgrade: %s", err.Error())
	}
	defer c.Close()

	mt, _, err := c.ReadMessage()
	if err != nil {
		log.Println("read: ", err)
		return
	}

	containers := handlers.List()

	format := "%-15s\t%-10s\t%s"
	msg := fmt.Sprintf(format, "Function ID", "State", "Service URL")
	err = c.WriteMessage(mt, []byte(msg))
	if err != nil {
		log.Println("write: ", err)
	}

	for _, container := range containers {
		var serviceURL string
		if len(container.Ports) > 0 {
			serviceURL = fmt.Sprintf("%s:%d", container.Ports[0].IP, container.Ports[0].PublicPort)
		}
		msg = fmt.Sprintf(format, container.ID[:10], container.State, serviceURL)
		err = c.WriteMessage(mt, []byte(msg))
		if err != nil {
			log.Println("write: ", err)
		}
	}

	closeConn(c, "0")
}

func down(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade: %s", err.Error())
	}
	defer c.Close()

	mt, message, err := c.ReadMessage()
	if err != nil {
		log.Println("read: ", err)
	}

	var ids []string
	if msg := string(message); msg != "*" {
		ids = strings.Split(msg, " ")
	}
	containers := handlers.List(ids...)

	count := len(containers)
	downResultCh := make(chan Message.DownMsgMeta, count)
	for _, container := range containers {
		ids = append(ids, container.ID)
		go handlers.Down(container.ID[:10], container.Image, downResultCh)
	}

	// collect down result
	var downs []Message.DownMsgMeta
	for downResult := range downResultCh {
		downs = append(downs, downResult)
		if len(downs) == count {
			close(downResultCh)
		}
	}

	msg := Message.CreateDownMessage(downs)
	c.WriteMessage(mt, []byte(msg))
	closeConn(c, "0")
}

func closeConn(c *websocket.Conn, msg string) {
	byteMsg := websocket.FormatCloseMessage(1000, msg)
	c.WriteMessage(websocket.CloseMessage, byteMsg)
}

// Start parses input and launches the fx server in a blocking process
func Start(verbose bool) {
	flag.Parse()
	log.SetFlags(0)

	env.Init(verbose)

	http.HandleFunc("/health", health)
	http.HandleFunc("/up", up)
	http.HandleFunc("/down", down)
	http.HandleFunc("/list", list)

	addr := fmt.Sprintf("%s:%s", config.Server["host"], config.Server["port"])
	log.Printf("fx serves on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
