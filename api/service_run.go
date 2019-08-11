package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/metrue/fx/types"
	"github.com/pkg/errors"
)

// ContainerCreateRequestPayload request paylaod
type ContainerCreateRequestPayload struct {
	*container.Config
	HostConfig       *container.HostConfig
	NetworkingConfig *network.NetworkingConfig
}

// Run a service
func (api *API) Run(port int, service *types.Service) error {
	config := &container.Config{
		Image:  service.Image,
		Labels: map[string]string{},
		ExposedPorts: nat.PortSet{
			"3000/tcp": struct{}{},
		},
	}
	hostConfig := &container.HostConfig{
		AutoRemove: true,
		PortBindings: nat.PortMap{
			"3000/tcp": []nat.PortBinding{
				{
					HostIP:   types.DefaultHost,
					HostPort: fmt.Sprintf("%d", port),
				},
			},
		},
	}
	req := ContainerCreateRequestPayload{
		Config:     config,
		HostConfig: hostConfig,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return errors.Wrap(err, "error mashal container create req")
	}

	// create container
	path := fmt.Sprintf("/containers/create?name=%s", service.Name)
	var createRes container.ContainerCreateCreatedBody
	if err := api.post(path, body, 201, &createRes); err != nil {
		return errors.Wrap(err, "create container request failed")
	}

	if createRes.ID == "" {
		return fmt.Errorf("container id is missing")
	}

	// start container
	path = fmt.Sprintf("/containers/%s/start", createRes.ID)
	url := fmt.Sprintf("%s%s", api.endpoint, path)
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return errors.Wrap(err, "error new container create request")
	}
	client := &http.Client{Timeout: 20 * time.Second}
	if _, err = client.Do(request); err != nil {
		return errors.Wrap(err, "error do start container request")
	}

	info, err := api.inspect(createRes.ID)
	if err != nil {
		msg := fmt.Sprintf("inspect container %s error", createRes.ID)
		return errors.Wrap(err, msg)
	}
	service.ID = info.ID
	service.Host = info.HostConfig.PortBindings["3000/tcp"][0].HostIP
	service.Port = port
	service.State = info.State.Status

	return nil
}
