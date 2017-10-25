package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"

	"./worker"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{} // use default options

func deploy() {
	log.Print("deployed")
}

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

	worker.Work(lang, body, c, mt)

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

	for {
		mt, _, err := c.ReadMessage()
		if err != nil {
			log.Println("read: ", err)
			break
		}
		worker.List(c, mt)
	}
}

func down(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("upgrade: ", err)
	}
	defer c.Close()

	mt, message, err := c.ReadMessage()
	if err != nil {
		log.Println("read: ", err)
	}
	ids := strings.Split(string(message), " ")

	doneCh := make(chan bool)
	msgCh := make(chan string)

	if ids[0] == "*" {
		fmt.Println("end all")
		go worker.StopAll(c, msgCh, doneCh)
		done := <-doneCh
		if done {
			c.WriteMessage(mt, []byte("All down"))
			closeConn(c, "0")
		} else {
			c.WriteMessage(mt, []byte("Could not down all"))
			closeConn(c, "0")
		}
	} else {
		fmt.Println("end list")
		for _, id := range ids {
			go worker.Stop(c, id, msgCh, doneCh, true)
		}
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

func main() {
	flag.Parse()
	log.SetFlags(0)

	http.HandleFunc("/health", health)
	http.HandleFunc("/up", up)
	http.HandleFunc("/down", down)
	http.HandleFunc("/list", list)

	log.Fatal(http.ListenAndServe(*addr, nil))

	log.Printf("addr: %p", *addr)
}
