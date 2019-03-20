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
	"github.com/phayes/freeport"
)

type healtCheck struct {
	Test        []string `json:"Test"`
	Interval    float64  `json:"Interval"`
	Timeout     float64  `json:"Timeout"`
	Retries     int64    `json:"Retries"`
	StartPeriod float64  `json:"StartPeriod"`
}

// ContainerCreateRequestPayload request paylaod
type ContainerCreateRequestPayload struct {
	Hostname         string                   `json:"Hostname,omitempty"`
	Domainname       string                   `json:"Domainname,omitempty"`
	User             string                   `json:"User,omitempty"`
	AttachStdin      bool                     `json:"AttachStdin,omitempty"`
	AttachStdout     bool                     `json:"AttachStdout,omitempty"`
	AttachStderr     bool                     `json:"AttachStderr,omitempty"`
	Tty              bool                     `json:"Tty,omitempty"`
	OpenStdin        bool                     `json:"OpenStdin,omitempty"`
	StdinOnce        bool                     `json:"StdinOnce,omitempty"`
	Env              []string                 `json:"Env,omitempty"`
	Cmd              []string                 `json:"Cmd,omitempty"`
	Entrypoint       string                   `json:"Entrypoint,omitempty"`
	Image            string                   `json:"Image,omitempty"`
	Labels           map[string]string        `json:"Labels,omitempty"`
	Volumes          map[string]interface{}   `json:"Volumes,omitempty"`
	Healthcheck      healtCheck               `json:"Healthcheck,omitempty"`
	WorkingDir       string                   `json:"WorkingDir,omitempty"`
	NetworkDisabled  bool                     `json:"NetworkDisabled,omitempty"`
	MacAddress       string                   `json:"MacAddress,omitempty"`
	ExposedPorts     nat.PortSet              `json:"ExposedPorts,omitempty"`
	StopSignal       string                   `json:"StopSignal,omitempty"`
	HostConfig       container.HostConfig     `json:"HostConfig,omitempty"`
	NetworkingConfig network.NetworkingConfig `json:"NetworkingConfig,omitempty"`
}

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
