package server

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	Config "../config"
	"./handlers"
	"./env"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

func health(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "I am OK, %s!", r.URL.Path[1:])
}

func up(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade: ", err)
	}

	log.Println("to up")
	defer c.Close()

	_, lang, err := c.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}

	mt, body, err := c.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}

	handlers.Up(lang, body, c, mt)

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Println("read:", msg)
	}
}

func list(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade: ", err)
	}
	defer c.Close()

	mt, _, err := c.ReadMessage()
	if err != nil {
		log.Println("read: ", err)
		return
	}

	containers := handlers.List()

	msg := "Function ID" + "\t" + "Service URL"
	err = c.WriteMessage(mt, []byte(msg))
	if err != nil {
		log.Println("write: ", err)
	}

	for _, container := range containers {
		msg = fmt.Sprintf("%s\t%s:%d", container.ID[:10], container.Ports[0].IP, container.Ports[0].PublicPort)
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
		log.Printf("upgrade: ", err)
	}
	defer c.Close()

	doneCh := make(chan bool)
	msgCh := make(chan string)

	mt, message, err := c.ReadMessage()
	if err != nil {
		log.Println("read: ", err)
	}

	var ids []string
	if string(message) == "*" {
		fmt.Println("end all")
		containers := handlers.List()
		ids = make([]string, len(containers))
		for i, container := range containers {
			ids[i] = container.ID[:10]
		}
	} else {
		fmt.Println("end list")
		ids = strings.Split(string(message), " ")
	}

	for _, id := range ids {
		go handlers.Down(id, msgCh, doneCh)
	}

	numSuccess := 0
	numFail := 0
	for {
		select {
		case newDone := <-doneCh:
			if newDone {
				numSuccess++
			} else {
				numFail++
			}

			if numSuccess+numFail == len(ids) {
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

func Start() {
	flag.Parse()
	log.SetFlags(0)

	env.Init();

	http.HandleFunc("/health", health)
	http.HandleFunc("/up", up)
	http.HandleFunc("/down", down)
	http.HandleFunc("/list", list)

	log.Printf("fx serves on %s", *Config.ServerAddr)
	log.Fatal(http.ListenAndServe(*Config.ServerAddr, nil))

	log.Printf("addr: %p", *Config.ServerAddr)
}
