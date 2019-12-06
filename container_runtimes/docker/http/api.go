package api

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/apex/log"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	dockerTypesContainer "github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
	containerruntimes "github.com/metrue/fx/container_runtimes"
	"github.com/metrue/fx/types"
	"github.com/metrue/fx/utils"
	"github.com/pkg/errors"
)

// API interact with dockerd http api
type API struct {
	endpoint string
	version  string
}

// Create a API
func Create(host string, port string) (*API, error) {
	addr := host + ":" + port
	v, err := version(addr)
	if err != nil {
		return nil, err
	}
	endpoint := fmt.Sprintf("http://%s:%s/v%s", host, port, v)
	return &API{
		endpoint: endpoint,
		version:  v,
	}, nil
}

// MustCreate a api object, panic if not
func MustCreate(host string, port string) *API {
	addr := host + ":" + port
	v, err := version(addr)
	if err != nil {
		panic(err)
	}
	endpoint := fmt.Sprintf("http://%s:%s/v%s", host, port, v)
	return &API{
		endpoint: endpoint,
		version:  v,
	}
}

func (api *API) get(path string, qs string, v interface{}) error {
	url := fmt.Sprintf("%s%s", api.endpoint, path)
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}
	if qs != "" {
		url += "?" + qs
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("request %s failed: %d - %s", url, resp.StatusCode, resp.Status)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &v)
	if err != nil {
		return err
	}
	return nil
}

func (api *API) post(path string, body []byte, expectStatus int, v interface{}) error {
	url := fmt.Sprintf("%s%s", api.endpoint, path)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != expectStatus {
		return fmt.Errorf("request %s (%s) failed: %d - %s", url, string(body), resp.StatusCode, resp.Status)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	return nil
}

// Version get version of docker engine
func (api *API) Version(ctx context.Context) (string, error) {
	return version(api.endpoint)
}

func version(endpoint string) (string, error) {
	path := endpoint + "/version"
	if !strings.HasPrefix(path, "http") {
		path = "http://" + path
	}

	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("request %s failed: %d - %s", path, resp.StatusCode, resp.Status)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var res dockerTypes.Version
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", err
	}
	return res.APIVersion, nil
}

// ListContainer list service
func (api *API) ListContainer(ctx context.Context, name string) ([]types.Service, error) {
	if name != "" {
		info, err := api.inspect(name)
		if err != nil {
			return []types.Service{}, err
		}

		port, err := strconv.Atoi(info.HostConfig.PortBindings["3000/tcp"][0].HostPort)
		if err != nil {
			return []types.Service{}, err
		}
		return []types.Service{
			types.Service{
				Name:  name,
				Image: info.Image,
				State: info.State.Status,
				ID:    info.ID,
				Host:  info.HostConfig.PortBindings["3000/tcp"][0].HostIP,
				Port:  port,
			},
		}, nil
	}

	type filterItem struct {
		Status []string `json:"status,omitempty"`
		Label  []string `json:"label,omitempty"`
		Name   []string `json:"name,omitempty"`
	}

	type Filters struct {
		Items string `url:"filters"`
	}

	filter := filterItem{
		// Status: []string{"running"},
		Label: []string{"belong-to=fx"},
	}

	q, err := json.Marshal(filter)
	if err != nil {
		return []types.Service{}, err
	}

	filters := Filters{Items: string(q)}
	qs, err := query.Values(filters)
	if err != nil {
		return []types.Service{}, err
	}

	var containers []dockerTypes.Container
	if err := api.get("/containers/json", qs.Encode(), &containers); err != nil {
		return []types.Service{}, err
	}

	svs := make(map[string]types.Service)
	for _, container := range containers {
		// container name have extra forward slash
		// https://github.com/moby/moby/issues/6705
		if strings.HasPrefix(container.Names[0], fmt.Sprintf("/%s", name)) {
			svs[container.Image] = types.Service{
				Name:  container.Names[0],
				Image: container.Image,
				ID:    container.ID,
				Host:  container.Ports[0].IP,
				Port:  int(container.Ports[0].PublicPort),
				State: container.State,
			}
		}
	}
	services := []types.Service{}
	for _, s := range svs {
		services = append(services, s)
	}

	return services, nil
}

// BuildImage build image
func (api *API) BuildImage(ctx context.Context, workdir string, name string) error {
	tarDir, err := ioutil.TempDir("/tmp", "fx-tar")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tarDir)

	imageID := uuid.New().String()
	tarFilePath := filepath.Join(tarDir, fmt.Sprintf("%s.tar", imageID))

	if err := utils.TarDir(workdir, tarFilePath); err != nil {
		return err
	}

	dockerBuildContext, err := os.Open(tarFilePath)
	if err != nil {
		return err
	}
	defer dockerBuildContext.Close()

	type buildQuery struct {
		Labels     string `url:"labels,omitempty"`
		Tags       string `url:"t,omitempty"`
		Dockerfile string `url:"dockerfile,omitempty"`
	}

	// Apply default labels
	labelsJSON, _ := json.Marshal(map[string]string{
		"belong-to": "fx",
	})
	q := buildQuery{
		Labels:     string(labelsJSON),
		Dockerfile: "Dockerfile",
	}

	qs, err := query.Values(q)
	if err != nil {
		return err
	}
	qs.Add("t", name)
	qs.Add("t", imageID)

	path := "/build"
	url := fmt.Sprintf("%s%s?%s", api.endpoint, path, qs.Encode())
	req, err := http.NewRequest("POST", url, dockerBuildContext)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-tar")
	client := &http.Client{Timeout: 600 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		if os.Getenv("DEBUG") != "" {
			log.Infof(scanner.Text())
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

// PushImage push a image
func (api *API) PushImage(ctx context.Context, name string) (string, error) {
	return "", nil
}

// InspectImage inspect image
func (api *API) InspectImage(ctx context.Context, name string, image interface{}) error {
	return nil
}

// TagImage tag image
func (api *API) TagImage(ctx context.Context, name string, tag string) error {
	query := url.Values{}
	query.Set("repo", name)
	query.Set("tag", tag)
	path := fmt.Sprintf("/images/%s/tag?%s", name, query.Encode())

	url := fmt.Sprintf("%s%s", api.endpoint, path)
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 10 * time.Second}
	if _, err = client.Do(req); err != nil {
		return err
	}

	return nil
}

// StartContainer start container
func (api *API) StartContainer(ctx context.Context, name string, image string, bindings []types.PortBinding) error {
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

	endpoint := &network.EndpointSettings{
		NetworkID: networks[0].ID,
	}
	networkConfig := &network.NetworkingConfig{
		EndpointsConfig: map[string]*network.EndpointSettings{
			"fx-net": endpoint,
		},
	}

	portSet := nat.PortSet{}
	portMap := nat.PortMap{}
	for _, binding := range bindings {
		bindings := []nat.PortBinding{
			nat.PortBinding{
				HostIP:   types.DefaultHost,
				HostPort: fmt.Sprintf("%d", binding.ServiceBindingPort),
			},
		}
		port := nat.Port(fmt.Sprintf("%d/tcp", binding.ContainerExposePort))
		portSet[port] = struct{}{}
		portMap[port] = bindings
	}
	config := &dockerTypesContainer.Config{
		Image:        image,
		ExposedPorts: portSet,
	}

	hostConfig := &dockerTypesContainer.HostConfig{
		AutoRemove:   true,
		PortBindings: portMap,
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
	path := fmt.Sprintf("/containers/create?name=%s", name)
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

	if _, err = api.inspect(createRes.ID); err != nil {
		msg := fmt.Sprintf("inspect container %s error", name)
		return errors.Wrap(err, msg)
	}

	return nil
}

// StopContainer stop a container
func (api *API) StopContainer(ctx context.Context, name string) error {
	return api.Stop(name)
}

// InspectContainer inspect container
func (api *API) InspectContainer(ctx context.Context, name string, container interface{}) error {
	path := fmt.Sprintf("/containers/%s/json", name)
	if err := api.get(path, "", &container); err != nil {
		return err
	}
	return nil
}

var (
	_ containerruntimes.ContainerRuntime = &API{}
)
