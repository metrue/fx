package worker

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	"github.com/phayes/freeport"
	"github.com/rs/xid"
)

func runCmd(cmdName string, cmdArgs []string) {
	cmd := exec.Command(cmdName, cmdArgs...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe for Cmd", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s | %s\n", cmdName, scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting Cmd", err)
		os.Exit(1)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for Cmd", err)
		os.Exit(1)
	}
}

func initWorkDirectory(dir string) {
	cmd := "./bin/init_worker_dir.sh"
	args := []string{dir}
	runCmd(cmd, args)
	fmt.Println("work dir initialized")
}

func dispatchFuncion(data []byte, dir string) {
	f, err := os.Create(dir + "/function/index.js")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	n, err := f.Write(data)
	if err != nil {
		panic(err)
	}

	log.Println("func recved: %s", n)
}

func buildService(name string, dir string) {
	cmd := "./bin/build_service.sh"
	args := []string{dir, name}
	runCmd(cmd, args)
	fmt.Println("service built")
}

func deployService(name string, dir string, port string) {
	cmd := "./bin/start_service.sh"
	args := []string{dir, name, port}
	runCmd(cmd, args)
	fmt.Println("service deployed")
}

func notify(connection *websocket.Conn, messageType int, message string) {
	err := connection.WriteMessage(messageType, []byte(message))
	if err != nil {
		log.Println("write: ", err)
	}
}

func closeConnection(connection *websocket.Conn) {
	connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Close"))
}

func Work(msg []byte, connection *websocket.Conn, messageType int) {
	go func() {
		var guid = xid.New().String()
		var dir = guid
		var name = guid

		port, err := freeport.GetFreePort()
		if err != nil {
			panic(err)
		}

		initWorkDirectory(dir)
		notify(connection, messageType, "work dir initialized")
		dispatchFuncion(msg, dir)
		notify(connection, messageType, "function dispatched")
		buildService(name, dir)
		notify(connection, messageType, "function built")
		deployService(name, dir, strconv.Itoa(port))
		notify(connection, messageType, "function deployed: http://localhost:"+strconv.Itoa(port))

		closeConnection(connection)
	}()
}

func List(connection *websocket.Conn, messageType int) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		msg := container.ID[:10] + " " + container.Ports[0].IP + strconv.Itoa(int(container.Ports[0].PublicPort))
		notify(connection, messageType, msg)
	}

	closeConnection(connection)
}
