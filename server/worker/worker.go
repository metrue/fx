package worker

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"time"

	"../utils"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/gorilla/websocket"
	"github.com/jhoonb/archivex"
	"github.com/phayes/freeport"
	"github.com/rs/xid"
)

var funcNames = map[string]string{
	"go":     "/fx.go",
	"node":   "/fx.js",
	"ruby":   "/fx.rb",
	"python": "/fx.py",
}

func initWorkDirectory(lang string, dir string) {
	err := utils.CopyDir(path.Join(os.Getenv("HOME"), ".fx/images", lang), dir)
	if err != nil {
		panic(err)
	}
}

func dispatchFuncion(lang string, data []byte, dir string) {
	fileName := funcNames[lang]
	f, err := os.Create(dir + fileName)
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

func checkerror(err error) {
	if err != nil {
		panic(err)
	}
}

func buildService(name string, dir string) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	tar := new(archivex.TarFile)
	tar.Create(dir)
	tar.AddAll(dir, false)
	tar.Close()
	dockerBuildContext, buildContextErr := os.Open(dir + ".tar")
	if buildContextErr != nil {
		panic(buildContextErr)
	}
	defer dockerBuildContext.Close()

	buildOptions := types.ImageBuildOptions{
		Dockerfile: "Dockerfile", // optional, is the default
		Tags:       []string{name},
	}
	buildResponse, buildErr := cli.ImageBuild(context.Background(), dockerBuildContext, buildOptions)
	if buildErr != nil {
		panic(buildErr)
	}
	log.Println("build ", buildResponse, buildErr)

	response, err := ioutil.ReadAll(buildResponse.Body)
	if err != nil {
		fmt.Printf("%s", err.Error())
	}
	fmt.Println(string(response))
}

func deployService(name string, dir string, port string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	imageName := name
	containerConfig := &container.Config{
		Image: imageName,
		ExposedPorts: nat.PortSet{
			"3000/tcp": struct{}{},
		},
	}
	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			"3000/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: port,
				},
			},
		},
	}

	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConfig, nil, "")
	if err != nil {
		panic(err)
	}
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
	fmt.Println(resp.ID)
}

func notify(connection *websocket.Conn, messageType int, message string) {
	err := connection.WriteMessage(messageType, []byte(message))
	if err != nil {
		log.Println("write: ", err)
	}
}

func closeConnection(connection *websocket.Conn) {
	connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "0"))
}

func Up(
	lang []byte,
	body []byte,
	connection *websocket.Conn,
	messageType int,
) {
	// go func() {
	var guid = xid.New().String()
	var dir = guid
	var name = guid

	port, err := freeport.GetFreePort()
	if err != nil {
		panic(err)
	}

	initWorkDirectory(string(lang), dir)
	notify(connection, messageType, "work dir initialized")
	dispatchFuncion(string(lang), body, dir)
	notify(connection, messageType, "function dispatched")
	buildService(name, dir)
	notify(connection, messageType, "function built")
	deployService(name, dir, strconv.Itoa(port))
	notify(connection, messageType, "function deployed: http://localhost:"+strconv.Itoa(port))

	closeConnection(connection)
	// }()
}

func List() []types.Container {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	return containers
}

func Stop(
	containID string,
	msgChan chan<- string,
	doneChan chan<- bool,
) {
	checkErr := func(err error) bool {
		if err != nil {
			log.Println(err)
			doneChan <- false
			return true
		}
		return false
	}

	cli, err := client.NewEnvClient()
	if checkErr(err) {
		return
	}

	timeout := time.Duration(1) * time.Second
	err = cli.ContainerStop(context.Background(), containID, &timeout)
	if checkErr(err) {
		return
	}

	fmt.Println("I am closed " + containID)
	msgChan <- containID + " Stopped"
	doneChan <- true
}
