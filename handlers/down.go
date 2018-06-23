package handlers

import (
	"github.com/metrue/fx/api"
	"github.com/metrue/fx/pkg/docker"
	"github.com/pkg/errors"
)

var (
	RemoveContainerError = errors.New("Failed to remove container")
	RemoveImageError     = errors.New("Failed to remove image")
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
		return res, RemoveContainerError
	}

	res.ContainerStatus = "stopped"
	err = docker.ImageRemove(image)
	if err != nil {
		return res, RemoveImageError
	}

	return res, nil
}
