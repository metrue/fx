package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/go-connections/nat"
	"github.com/gobuffalo/packr"
	"github.com/google/go-querystring/query"
	"github.com/google/uuid"
	"github.com/metrue/fx/types"
	"github.com/metrue/fx/utils"
	"github.com/phayes/freeport"
)

const serviceNamePrefix = "fx_"

// API interact with dockerd http api
type API struct {
	endpoint string
	box      packr.Box
}

// NewWithDockerRemoteAPI create a api with docker remote api
func NewWithDockerRemoteAPI(url string, box packr.Box) *API {
	return &API{endpoint: url, box: box}
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
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}

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

// Version get version of dockerd server
func (api *API) version() (string, error) {
	var res dockerTypes.Version
	if err := api.get("/version", "", &res); err != nil {
		return "", err
	}
	return res.APIVersion, nil
}

// Build build a project
func (api *API) Build(project types.Project) (types.Service, error) {
	dir, err := ioutil.TempDir("/tmp", "fx-build-dir")
	if err != nil {
		return types.Service{}, err
	}

	defer os.RemoveAll(dir)

	for _, file := range project.Files {
		tmpfn := filepath.Join(dir, file.Path)
		if err := utils.EnsureFile(tmpfn); err != nil {
			return types.Service{}, err
		}
		if err := ioutil.WriteFile(tmpfn, []byte(file.Body), 0666); err != nil {
			return types.Service{}, err
		}
	}

	tarDir, err := ioutil.TempDir("/tmp", "fx-tar")
	if err != nil {
		return types.Service{}, err
	}

	imageID := uuid.New().String()
	tarFilePath := filepath.Join(tarDir, fmt.Sprintf("%s.tar", imageID))
	if err := utils.TarDir(dir, tarFilePath); err != nil {
		return types.Service{}, err
	}

	dockerBuildContext, err := os.Open(tarFilePath)
	if err != nil {
		return types.Service{}, err
	}
	defer dockerBuildContext.Close()

	type buildQuery struct {
		Labels     string `url:"labels"`
		Tags       string `url:"t"`
		Dockerfile string `url:"dockerfile"`
	}

	// Apply default labels
	labelsJSON, _ := json.Marshal(
		map[string]string{
			"belong-to": "fx",
		},
	)

	q := buildQuery{
		Tags:       imageID,
		Labels:     string(labelsJSON),
		Dockerfile: "Dockerfile",
	}
	qs, err := query.Values(q)
	if err != nil {
		return types.Service{}, err
	}

	version, err := api.version()
	if err != nil {
		return types.Service{}, err
	}
	path := "/build"
	url := fmt.Sprintf("http://%s/v%s%s?%s", api.endpoint, version, path, qs.Encode())
	req, err := http.NewRequest("POST", url, dockerBuildContext)
	if err != nil {
		return types.Service{}, err
	}

	req.Header.Set("Content-Type", "application/x-tar")
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return types.Service{}, err
	}

	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		// TODO Maybe need log something out
		// fmt.Printf("%s\n", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return types.Service{}, err
	}

	return types.Service{
		Name:   project.Name,
		Status: types.ServiceStatusInit,
		Image:  imageID,
	}, nil
}

// List list service
func (api *API) list(name string) ([]types.Service, error) {
	version, err := api.version()
	if err != nil {
		return []types.Service{}, err
	}

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
				Name:   name,
				Image:  info.Image,
				Status: types.ServiceStatusRunning,
				Instances: []types.Instance{
					types.Instance{
						ID:    info.ID,
						Host:  info.HostConfig.PortBindings["3000/tcp"][0].HostIP,
						Port:  port,
						State: info.State.Status,
					},
				},
			},
		}, nil
	}

	type filterItem struct {
		Status []string `json:"url,omitempty"`
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
	path := fmt.Sprintf("/v%s/containers/json", version)
	if err := api.get(path, qs.Encode(), &containers); err != nil {
		return []types.Service{}, err
	}

	svs := make(map[string]types.Service)
	for _, container := range containers {
		// container name have extra forward slash
		// https://github.com/moby/moby/issues/6705
		if strings.HasPrefix(container.Names[0], fmt.Sprintf("/%s", name)) {
			instance := types.Instance{
				ID:    container.ID,
				Host:  container.Ports[0].IP,
				Port:  int(container.Ports[0].PublicPort),
				State: container.State,
			}
			if svs[container.Image].Instances != nil {
				instances := append(svs[container.Image].Instances, instance)
				svs[container.Image] = types.Service{Instances: instances}
			} else {
				svs[container.Image] = types.Service{
					Name:      name,
					Image:     container.Image,
					Status:    types.ServiceStatusRunning,
					Instances: []types.Instance{instance},
				}
			}
		}
	}
	services := []types.Service{}
	for _, s := range svs {
		services = append(services, s)
	}

	return services, nil
}

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

	version, err := api.version()
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/v%s/containers/create?name=%s", version, service.Name)
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

	path = fmt.Sprintf("/v%s/containers/%s/start", version, res.ID)
	url := fmt.Sprintf("http://%s%s", api.endpoint, path)
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

// Stop a container by name
func (api *API) Stop(name string) error {
	version, err := api.version()
	if err != nil {
		return err
	}
	path := fmt.Sprintf("/v%s/containers/%s/stop", version, name)
	url := fmt.Sprintf("http://%s%s", api.endpoint, path)
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{Timeout: 20 * time.Second}
	_, err = client.Do(request)
	if err != nil {
		return err
	}

	return nil
}
