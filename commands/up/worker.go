package up

import (
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/metrue/fx/log"
)

// Worker handles a functional service
type Worker struct {
	conn   *websocket.Conn
	ch     chan<- bool
	src    string
	lang   string
	logger *Logger
	dead   bool
}

// NewWorker creates and returns a new worker
func NewWorker(src, lang string, conn *websocket.Conn, ch chan<- bool) *Worker {
	worker := &Worker{
		src:    src,
		lang:   lang,
		conn:   conn,
		ch:     ch,
		logger: log.NewLogger("[" + src + "]"),
	}
	worker.conn.SetCloseHandler(worker.closeHandler)
	return worker
}

func (worker *Worker) checkErr(err error) bool {
	if err != nil {
		worker.logger.Err(err)
		if !websocket.IsCloseError(err, 1000) {
			worker.ch <- false
		}
		return true
	}
	return false
}

func (worker *Worker) closeHandler(
	code int,
	msg string,
) error {
	worker.dead = true
	if msg == "0" {
		worker.ch <- true
		worker.logger.Log(msg)
	} else {
		worker.ch <- false
		worker.logger.Err(msg)
	}
	return nil
}

// Work starts and handles a function from Worker's information
func (worker *Worker) Work() {
	if worker.dead {
		return
	}
	logger := worker.logger
	conn := worker.conn
	logger.Log("Deploying...")
	// Open function source file
	file, err := os.Open(worker.src)
	if worker.checkErr(err) {
		return
	}
	defer file.Close()
	// Send source language type
	err = conn.WriteMessage(
		websocket.TextMessage,
		[]byte(worker.lang),
	)
	if worker.checkErr(err) {
		return
	}

	// Get websocket connection writer
	writer, err := conn.NextWriter(
		websocket.TextMessage,
	)
	if worker.checkErr(err) {
		return
	}

	// Send function source file
	bytesSent, err := io.Copy(
		writer,
		file,
	)
	if worker.checkErr(err) {
		return
	}
	logger.Log(
		fmt.Sprintf(
			"Sent bytes: %d",
			bytesSent,
		))
	err = writer.Close()
	if worker.checkErr(err) {
		return
	}

	interrupt := make(chan os.Signal, 2)
	signal.Notify(interrupt, os.Interrupt, os.Kill)

	errChan := make(chan error)
	msgChan := make(chan string)
	readReply := func(
		c *websocket.Conn,
		msgCh chan<- string,
		errCh chan<- error,
	) {
		_, msg, err := c.ReadMessage()
		if err != nil {
			errCh <- err
			return
		}
		msgCh <- string(msg)
	}
	go readReply(conn, msgChan, errChan)

	// Wait for deploy done
	for {
		select {
		case <-interrupt:
			closeMsg := websocket.FormatCloseMessage(1000, "1")
			conn.WriteMessage(websocket.CloseMessage, closeMsg)
		case newMsg := <-msgChan:
			logger.Log(newMsg)
			if worker.dead {
				return
			}
			go readReply(conn, msgChan, errChan)
		case newErr := <-errChan:
			if worker.checkErr(newErr) {
				return
			}
		}
	}
}
