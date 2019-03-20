package api

import (
	"fmt"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

type containerInfo struct {
	ID         string                     `json:"Id"`
	State      dockerTypes.ContainerState `json:"State"`
	Image      string                     `json:"Image"`
	HostConfig container.HostConfig       `json:"HostConfig"`
}

func (api *API) inspect(identify string) (containerInfo, error) {
	var info containerInfo

	path := fmt.Sprintf("/containers/%s/json", identify)
	if err := api.get(path, "", &info); err != nil {
		return info, err
	}

	return info, nil
}
