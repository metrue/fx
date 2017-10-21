package uper

import (
	"fmt"
	"os"
	"io"
	"os/signal"

	"../logger"

	"github.com/gorilla/websocket"
)

type Uper struct {
	conn *websocket.Conn
	ch chan<- bool
	src string
	logger *logger.Logger
	dead bool
}

func New(
	src string,
	conn *websocket.Conn,
	ch chan<- bool,
) *Uper {
	uper := &Uper{
		dead: false,
		src: src,
		conn: conn,
		ch: ch,
		logger: logger.New("[" + src + "]"),
	}
	uper.conn.SetCloseHandler(uper.closeHandler)
	return uper
}

func (uper *Uper)	checkErr(err error) bool {
	if err != nil {
		uper.logger.Err(err)
		if !websocket.IsCloseError(err, 1000) {
			uper.ch <- false
		}
		return true
	}
	return false
}

func (uper *Uper) closeHandler(
	code int,
	msg string,
) error {
	uper.dead = true
	if msg == "0" {
		uper.ch <- true
		uper.logger.Log(msg)
	} else {
		uper.ch <- false
		uper.logger.Err(msg)
	}
	return nil
}

func (uper *Uper) Work() {
	logger := uper.logger
	conn := uper.conn
	logger.Log("Deploying...")

	// Open function source file
	if uper.dead { return }
	file, err := os.Open(uper.src)
	if uper.checkErr(err) { return }
	defer file.Close()

	// Get websocket connection writer
	if uper.dead { return }
	writer, err := conn.NextWriter(
		websocket.TextMessage,
	)
	if uper.checkErr(err) { return }

	// Send function source file
	if uper.dead { return }
	bytesSent, err := io.Copy(
		writer,
		file,
	)
	if uper.checkErr(err) { return }
	logger.Log(
		fmt.Sprintf(
			"Sent bytes: %d",
			bytesSent,
		))
	err = writer.Close()
	if uper.checkErr(err) { return }


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
			if uper.dead { return }
			go readReply(conn, msgChan, errChan)
		case newErr := <-errChan:
			if uper.checkErr(newErr) { return }
		}
	}
}
