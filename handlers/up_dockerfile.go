package handlers

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/docker-api"
	"github.com/metrue/fx/utils"
	"github.com/phayes/freeport"
	"github.com/rs/xid"
)

// Up spins up a new function
func UpDockerfile(dockerfile api.DockerfileMeta) (*api.UpMsgMeta, error) {
	port, err := freeport.GetFreePort()
	if err != nil {
		return nil, err
	}

	var guid = xid.New().String()
	var dir = path.Join(os.TempDir(), "fx-", guid)
	defer cleanup(dir)

	var name = guid
	targetPath := path.Join(dir, "Dockerfile")
	writeErr := ioutil.WriteFile(targetPath, []byte(dockerfile.Content), 0644)
	if writeErr != nil {
		return nil, err
	}

	err = docker.Build(name, dir)
	if err != nil {
		return nil, err
	}

	containerInfo, err := docker.Deploy(name, dir, strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	localAddr := fmt.Sprintf("127.0.0.1:%s", strconv.Itoa(port))
	remoteAddr := fmt.Sprintf("%s:%s", utils.GetHostIP().String(), strconv.Itoa(port))

	res := &api.UpMsgMeta{
		FunctionID:     containerInfo.ID[:10],
		FunctionSource: string("Dockerfile"),
		LocalAddress:   localAddr,
		RemoteAddress:  remoteAddr,
	}

	return res, nil
}
