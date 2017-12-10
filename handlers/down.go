package handlers

import (
	"github.com/metrue/fx/api"
	"github.com/metrue/fx/docker-api"
)

// Down stops the processes designated by a function
func Down(containerId string, image string, result chan<- api.DownMsgMeta) {
	res := api.DownMsgMeta{
		ContainerId:     containerId,
		ContainerStatus: "",
		ImageStatus:     "",
	}
	err := docker.Remove(containerId)
	if err == nil {
		res.ContainerStatus = "stopped"
	}

	if err := docker.ImageRemove(image); err != nil {
		res.ImageStatus = "not removed"
	} else {
		res.ImageStatus = "removed"
	}
	result <- res
}
