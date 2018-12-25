package service

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types"
	"github.com/metrue/fx/api"
	"github.com/metrue/fx/pkg/docker"
)

var (
	RemoveContainerError = errors.New("Failed to remove container")
	RemoveImageError     = errors.New("Failed to remove image")
)

// Down stops the processes designated by a function
func DoDown(containerID string, image string) (*api.DownMsgMeta, error) {
	res := &api.DownMsgMeta{
		ContainerId:     containerID,
		ContainerStatus: "",
		ImageStatus:     "removed",
	}

	err := docker.Remove(containerID)
	if err != nil {
		return res, RemoveContainerError
	}

	res.ContainerStatus = "stopped"
	err = docker.ImageRemove(image)
	if err != nil {
		return res, RemoveImageError
	}

	return res, nil
}

//downTask wrap a DownMsgMeta and an error from its processing
type downTask struct {
	Val *api.DownMsgMeta
	Err error
}

//newDownTask initialize a new downTask
func newDownTask(val *api.DownMsgMeta, err error) downTask {
	return downTask{
		Val: val,
		Err: err,
	}
}

//Down handle function removal requests
func Down(ctx context.Context, req *api.DownRequest) (*api.DownResponse, error) {

	// handle fx down *
	var ids []string
	if req.ID != nil && len(req.ID) > 0 {
		if req.ID[0] != "*" {
			ids = req.ID
		}
	}

	containers, err := DoList(ids...)
	if err != nil {
		return nil, err
	}
	count := len(containers)
	results := make(chan downTask, count)
	downResponse := &api.DownResponse{}

	if count == 0 {
		return downResponse, nil
	}

	for _, c := range containers {
		go func(container types.Container) {
			results <- newDownTask(DoDown(container.ID[:10], container.Image))
		}(c)
	}

	// collect down result
	var downs []*api.DownMsgMeta
	for result := range results {
		downResult := result.Val
		if result.Err != nil {
			downResult.Error = result.Err.Error()
		}
		downs = append(downs, downResult)
		if len(downs) == count {
			close(results)
		}
	}

	downResponse.Instances = downs
	return downResponse, nil
}
