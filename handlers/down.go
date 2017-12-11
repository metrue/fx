package handlers

import (
	"github.com/metrue/fx/api"
	"github.com/metrue/fx/docker-api"
)

// Down stops the processes designated by a function
func Down(containerID string, image string) (*api.DownMsgMeta, error) {

	res := &api.DownMsgMeta{
		ContainerId:     containerID,
		ContainerStatus: "",
		ImageStatus:     "",
	}

	err := docker.Remove(containerID)
	if err != nil {
		return nil, err
	}
	res.ContainerStatus = "stopped"

	err = docker.ImageRemove(image)
	if err != nil {
		return nil, err
	}
	res.ImageStatus = "removed"

	return res, err
}
