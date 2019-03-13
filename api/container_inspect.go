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

	version, err := api.version()
	if err != nil {
		return info, err
	}

	path := fmt.Sprintf("/v%s/containers/%s/json", version, identify)
	type containerInfo struct {
		ID         string                     `json:"Id"`
		State      dockerTypes.ContainerState `json:"State"`
		Image      string                     `json:"Image"`
		HostConfig container.HostConfig       `json:"HostConfig"`
	}
	if err := api.get(path, "", &info); err != nil {
		return info, err
	}

	return info, nil
}
