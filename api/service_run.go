package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/apex/log"
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

const fxNetworkName = "fx-net"

// Run a service
func (api *API) Run(port int, service *types.Service) error {
	networks, err := api.GetNetwork(fxNetworkName)
	if err != nil {
		return errors.Wrapf(err, "get network failed: %s", err)
	}

	if len(networks) == 0 {
		if err := api.CreateNetwork(fxNetworkName); err != nil {
			return errors.Wrapf(err, "error create network: %s", err)
		}
	}
	networks, _ = api.GetNetwork(fxNetworkName)
	config := &container.Config{
		Image: service.Image,
		ExposedPorts: nat.PortSet{
			"3000/tcp": struct{}{},
		},
	}

	endpoint := &network.EndpointSettings{
		NetworkID: networks[0].ID,
	}
	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"fx-net": endpoint,
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
		Config:           config,
		HostConfig:       hostConfig,
		NetworkingConfig: networkConfig,
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

	log.Infof("container %s created", service.Name)

	// start container
	path = fmt.Sprintf("/containers/%s/start", createRes.ID)
	url := fmt.Sprintf("%s%s", api.endpoint, path)
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return errors.Wrap(err, "error new container create request")
	}
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(request)
	if err != nil {
		return errors.Wrap(err, "error do start container request")
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if len(b) != 0 {
		msg := fmt.Sprintf("start container met issue: %s", string(b))
		return errors.New(msg)
	}
	log.Infof("container %s started", service.Name)

	info, err := api.inspect(createRes.ID)
	if err != nil {
		msg := fmt.Sprintf("inspect container %s error", service.Name)
		return errors.Wrap(err, msg)
	}
	service.ID = info.ID
	service.Host = info.HostConfig.PortBindings["3000/tcp"][0].HostIP
	service.Port = port
	service.State = info.State.Status

	return nil
}
