package handlers

import (
	"github.com/metrue/fx/docker-api"
	Message "github.com/metrue/fx/message"
)

// Down stops the processes designated by a function
func Down(containerId string, image string, result chan<- Message.DownMsgMeta) {
	res := Message.DownMsgMeta{
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
