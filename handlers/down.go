package handlers

import (
	"log"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/docker-api"
)

// Down stops the processes designated by a function
func Down(containerID string, image string) (*api.DownMsgMeta, error) {

	res := &api.DownMsgMeta{
		ContainerId:     containerID,
		ContainerStatus: "",
		ImageStatus:     "removed",
	}

	log.Printf("Removing container `%s`", containerID)
	err := docker.Remove(containerID)
	if err != nil {
		log.Printf("Failed to remove container `%s`: %s", containerID, err.Error())
		return res, err
	}
	res.ContainerStatus = "stopped"

	log.Printf("Removing image `%s`", image)
	err = docker.ImageRemove(image)
	if err != nil {
		log.Printf("Failed to remove image `%s`: %s", image, err.Error())
		return res, err
	}

	return res, err
}
