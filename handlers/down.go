package handlers

import (
	"github.com/metrue/fx/api"
	"github.com/metrue/fx/docker-api"
	"github.com/pkg/errors"
)

// Down stops the processes designated by a function
func Down(containerID string, image string) (*api.DownMsgMeta, error) {
	res := &api.DownMsgMeta{
		ContainerId:     containerID,
		ContainerStatus: "",
		ImageStatus:     "removed",
	}

	err := docker.Remove(containerID)
	if err != nil {
		err = errors.Wrap(err, "Failed to remove container")
		return res, err
	}

	res.ContainerStatus = "stopped"
	err = docker.ImageRemove(image)
	if err != nil {
		err = errors.Wrap(err, "[down.go] Failed to remove image")
		return res, err
	}

	return res, nil
}
