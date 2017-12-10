package handlers

import (
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/docker-api"
	"github.com/metrue/fx/image"
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
func Up(funcMeta api.FunctionMeta, result chan<- api.UpMsgMeta) {
	port, err := freeport.GetFreePort()
	if err != nil {
		panic(err)
	}
	var guid = xid.New().String()
	var dir = path.Join(os.TempDir(), "fx-", guid)
	defer cleanup(dir)
	var name = guid
	image.Get(dir, string(funcMeta.Lang), []byte(funcMeta.Content))
	docker.Build(name, dir)
	docker.Deploy(name, dir, strconv.Itoa(port))

	localAddr := fmt.Sprintf("127.0.0.1:%s", strconv.Itoa(port))
	remoteAddr := fmt.Sprintf("%s:%s", utils.GetHostIP().String(), strconv.Itoa(port))
	res := api.UpMsgMeta{
		FunctionSource: string(funcMeta.Path),
		LocalAddress:   localAddr,
		RemoteAddress:  remoteAddr,
	}
	result <- res
}
