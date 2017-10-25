package up

import (
	"fmt"
	"os"
	"io"
	"os/signal"

	"github.com/gorilla/websocket"
)

type Worker struct {
	conn *websocket.Conn
	ch chan<- bool
	src string
	lang string
	logger *Logger
	dead bool
}

func NewWorker(
	src string,
	lang string,
	conn *websocket.Conn,
	ch chan<- bool,
) *Worker {
	worker := &Worker{
		dead: false,
		src: src,
		lang: lang,
		conn: conn,
		ch: ch,
		logger: NewLogger("[" + src + "]"),
	}
	worker.conn.SetCloseHandler(worker.closeHandler)
	return worker
}

func (worker *Worker)	checkErr(err error) bool {
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

func (worker *Worker) Work() {
	logger := worker.logger
	conn := worker.conn
	logger.Log("Deploying...")

	// Open function source file
	if worker.dead { return }
	file, err := os.Open(worker.src)
	if worker.checkErr(err) { return }
	defer file.Close()

	// Send source language type
	if worker.dead { return }
	err = conn.WriteMessage(
		websocket.TextMessage,
		[]byte(worker.lang),
	)
	if worker.checkErr(err) { return }

	// Get websocket connection writer
	if worker.dead { return }
	writer, err := conn.NextWriter(
		websocket.TextMessage,
	)
	if worker.checkErr(err) { return }

	// Send function source file
	if worker.dead { return }
	bytesSent, err := io.Copy(
		writer,
		file,
	)
	if worker.checkErr(err) { return }
	logger.Log(
		fmt.Sprintf(
			"Sent bytes: %d",
			bytesSent,
		))
	err = writer.Close()
	if worker.checkErr(err) { return }


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
			if worker.dead { return }
			go readReply(conn, msgChan, errChan)
		case newErr := <-errChan:
			if worker.checkErr(newErr) { return }
		}
	}
}
