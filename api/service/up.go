package service

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/image"
	"github.com/metrue/fx/pkg/docker"
	"github.com/metrue/fx/pkg/utils"
	"github.com/phayes/freeport"
	"github.com/rs/xid"
)

//upTask wrap the UpMsgMeta and an error from processing
type upTask struct {
	Val *api.UpMsgMeta
	Err error
}

//newUpTask initialize a new upTask
func newUpTask(val *api.UpMsgMeta, err error) upTask {
	return upTask{
		Val: val,
		Err: err,
	}
}

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

// DoUp spins up a new function
func DoUp(funcMeta api.FunctionMeta) (*api.UpMsgMeta, error) {
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
		FunctionID:    containerInfo.ID[:10],
		LocalAddress:  localAddr,
		RemoteAddress: remoteAddr,
	}

	return res, nil
}

//Up deploy and run functions
func Up(ctx context.Context, req *api.UpRequest) (*api.UpResponse, error) {

	funcList := req.Functions

	count := len(funcList)
	results := make(chan upTask, count)

	for _, meta := range funcList {
		go func(funcMeta *api.FunctionMeta) {
			results <- newUpTask(DoUp(*funcMeta))
		}(meta)
	}

	// collect up results
	var ups []*api.UpMsgMeta
	for result := range results {
		upResult := result.Val
		if result.Err != nil {
			upResult = &api.UpMsgMeta{
				Error: result.Err.Error(),
			}
		}
		ups = append(ups, upResult)
		if len(ups) == count {
			close(results)
		}
	}

	return &api.UpResponse{
		Instances: ups,
	}, nil
}
