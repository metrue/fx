package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	"github.com/metrue/fx/types"
	"github.com/phayes/freeport"
)

// Run a service
func (api *API) Run(service *types.Service) error {
	port, err := freeport.GetFreePort()
	if err != nil {
		return err
	}

	req := ContainerCreateRequestPayload{
		Image:  service.Image,
		Labels: map[string]string{},
		ExposedPorts: nat.PortSet{
			"3000/tcp": struct{}{},
		},
		HostConfig: container.HostConfig{
			AutoRemove: true,
			PortBindings: nat.PortMap{
				"3000/tcp": []nat.PortBinding{
					{
						HostIP:   types.DefaultHost,
						HostPort: fmt.Sprintf("%d", port),
					},
				},
			},
		},
	}

	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	path := fmt.Sprintf("/containers/create?name=%s", service.Name)
	type containerCreateResponse struct {
		ID       string   `json:"Id"`
		Warnings []string `json:"Warnings"`
	}
	var res containerCreateResponse
	err = api.post(path, body, 201, &res)
	if err != nil {
		return err
	}

	if res.ID == "" {
		return fmt.Errorf("container id is missing")
	}

	path = fmt.Sprintf("/containers/%s/start", res.ID)
	url := fmt.Sprintf("%s%s", api.endpoint, path)
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 20 * time.Second}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	info, err := api.inspect(service.Name)
	if err != nil {
		return err
	}
	instance := types.Instance{
		ID:    info.ID,
		Host:  info.HostConfig.PortBindings["3000/tcp"][0].HostIP,
		Port:  port,
		State: info.State.Status,
	}
	service.Instances = append(service.Instances, instance)

	return nil
}
