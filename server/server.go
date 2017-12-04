package server

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/metrue/fx/config"
	"github.com/metrue/fx/server/env"
	"github.com/metrue/fx/server/handlers"

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

	_, lang, err := c.ReadMessage()
	if err != nil {
		log.Printf("read error: %s", err.Error())
		return
	}

	mt, body, err := c.ReadMessage()
	if err != nil {
		log.Printf("read error: %s", err.Error())
		return
	}

	handlers.Up(lang, body, c, mt)

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Printf("read error: %s", err.Error())
			return
		}
		log.Println("read:", msg)
	}
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

	doneCh := make(chan bool)
	msgCh := make(chan string)

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
	for _, container := range containers {
		ids = append(ids, container.ID)
		go handlers.Down(container.ID[:10], container.Image, msgCh, doneCh)
	}

	numSuccess := 0
	numFail := 0
	defer func() {
		close(doneCh)
		close(msgCh)
	}()
	for {
		select {
		case newDone := <-doneCh:
			if newDone {
				numSuccess++
			} else {
				numFail++
			}

			if numSuccess+numFail == count {
				res := fmt.Sprintf("Succed: %d", numSuccess)
				c.WriteMessage(mt, []byte(res))
				res = fmt.Sprintf("Failed: %d", numFail)
				c.WriteMessage(mt, []byte(res))
				closeConn(c, "0")
				return
			}
		case newMsg := <-msgCh:
			c.WriteMessage(mt, []byte(newMsg))
		}
	}
}

func closeConn(c *websocket.Conn, msg string) {
	byteMsg := websocket.FormatCloseMessage(1000, msg)
	c.WriteMessage(websocket.CloseMessage, byteMsg)
}

// Start parses input and launches the fx server in a blocking process
func Start() {
	flag.Parse()
	log.SetFlags(0)

	env.Init(true)

	http.HandleFunc("/health", health)
	http.HandleFunc("/up", up)
	http.HandleFunc("/down", down)
	http.HandleFunc("/list", list)

	addr := fmt.Sprintf("%s:%s", config.Server["host"], config.Server["port"])
	log.Printf("fx serves on %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
