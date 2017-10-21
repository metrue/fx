package worker

import (
	"archive/tar"
	"bufio"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/gorilla/websocket"
	"github.com/jhoonb/archivex"
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

func checkerror(err error) {
	if err != nil {
		panic(err)
	}
}

func tarDir(srcDir string, desFileName string) {
	dir, err := os.Open(srcDir)
	if err != nil {
		panic(err)
	}
	defer dir.Close()

	files, err := dir.Readdir(0)
	if err != nil {
		panic(err)
	}

	// create tar file
	tarfile, err := os.Create(desFileName)
	if err != nil {
		panic(err)
	}
	defer tarfile.Close()

	var fileWriter io.WriteCloser = tarfile

	tarfileWriter := tar.NewWriter(fileWriter)
	defer tarfileWriter.Close()

	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			continue
		}

		file, err := os.Open(dir.Name() + string(filepath.Separator) + fileInfo.Name())
		checkerror(err)
		defer file.Close()

		// prepare the tar header
		header := new(tar.Header)
		header.Name = file.Name()
		header.Size = fileInfo.Size()
		header.Mode = int64(fileInfo.Mode())
		header.ModTime = fileInfo.ModTime()

		err = tarfileWriter.WriteHeader(header)
		checkerror(err)

		_, err = io.Copy(tarfileWriter, file)
		checkerror(err)
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
	fmt.Println("try to deploy service")
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
	connection.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "0"))
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

func Stop(
	connection *websocket.Conn,
	containID string,
	msgChan chan<- string,
	done chan<- bool,
) {
	// notify(connection, websocket.TextMessage, "to stop"+containID)
	msgChan <- "to stop"+containID

	checkErr := func(err error) bool {
		if err != nil {
			log.Println(err)
			done <- false
			return true
			// panic(err)
		}
		return false
	}

	cli, err := client.NewEnvClient()
	if checkErr(err) { return }

	timeout := time.Duration(1) * time.Second
	// notify(connection, websocket.TextMessage, "to stop"+containID)
	err = cli.ContainerStop(context.Background(), containID, &timeout)
	if checkErr(err) { return }

	// notify(connection, websocket.TextMessage, "to stop"+containID)
	// msg := containID + " Stopped"
	msgChan <- containID + " Stopped"
	// notify(connection, websocket.TextMessage, msg)
	// closeConnection(connection)

	done <- true
}
