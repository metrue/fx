package handlers

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/image"
	"github.com/metrue/fx/pkg/docker"
	"github.com/metrue/fx/utils"
	"github.com/phayes/freeport"
	"github.com/rs/xid"
)

func cleanup(dir string) {
	format := "cleanup temp file [%s] error: %s\n"
	if err := os.RemoveAll(dir); err != nil {
		log.Printf(format, dir, err.Error())
	}
	dirTar := dir + ".tar"
	if err := os.RemoveAll(dirTar); err != nil {
		log.Printf(format, dirTar, err.Error())
	}
}

// Up spins up a new function
func Up(funcMeta api.FunctionMeta) (*api.UpMsgMeta, error) {

	port, err := freeport.GetFreePort()
	if err != nil {
		return nil, err
	}

	var guid = xid.New().String()
	var dir = path.Join(os.TempDir(), "fx-", guid)
	defer cleanup(dir)

	var name = guid
	err = image.Get(dir, string(funcMeta.Lang), []byte(funcMeta.Content))
	if err != nil {
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
		FunctionSource: string(funcMeta.Path),
		LocalAddress:   localAddr,
		RemoteAddress:  remoteAddr,
	}

	return res, nil
}
