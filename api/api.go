package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	dockerTypes "github.com/docker/docker/api/types"
	"github.com/gobuffalo/packr"
	"github.com/google/go-querystring/query"
	"github.com/metrue/fx/types"
)

const serviceNamePrefix = "fx_"

// API interact with dockerd http api
type API struct {
	endpoint string
	version  string
	box      packr.Box
}

// NewWithDockerRemoteAPI create a api with docker remote api
func NewWithDockerRemoteAPI(url string, version string) *API {
	box := packr.NewBox("./images")
	endpoint := fmt.Sprintf("%s/v%s", url, version)
	return &API{
		endpoint: endpoint,
		box:      box,
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

// List list service
func (api *API) list(name string) ([]types.Service, error) {
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
	if err := api.get("/containers/json", qs.Encode(), &containers); err != nil {
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
