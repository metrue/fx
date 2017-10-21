package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

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

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		worker.Work(message, c, mt)
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
	fmt.Fprint(w, "not ready yet", r.URL.Path[1:])
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
