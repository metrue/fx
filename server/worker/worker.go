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
	"path"
	"path/filepath"
	"strconv"
	"time"

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

func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

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

func initWorkDirectory(lang string, dir string) {
	// err := os.MkdirAll(dir, os.ModePerm)
	// if err != nil {
	// 	panic(err)
	// }

	scriptPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	err = CopyDir(path.Join(scriptPath, "..", "images", lang), dir)
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

func List(connection *websocket.Conn, messageType int) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	msg := "Function ID" + "\t" + "Service URL"
	notify(connection, websocket.TextMessage, msg)
	for _, container := range containers {
		msg := container.ID[:10] + "\t" + container.Ports[0].IP + ":" + strconv.Itoa(int(container.Ports[0].PublicPort))
		notify(connection, messageType, msg)
	}

	closeConnection(connection)
}

func StopAll(
	connection *websocket.Conn,
	msgChan chan<- string,
	done chan<- bool,
) {
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Println("end" + container.ID[:10])
		go Stop(connection, container.ID[:10], msgChan, done, false)
	}
	done <- true
}

func Stop(
	connection *websocket.Conn,
	containID string,
	msgChan chan<- string,
	done chan<- bool,
	sendDone bool,
) {
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
	if sendDone {
		done <- true
	} else {
		done <- false
	}
}
