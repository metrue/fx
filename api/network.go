package api

import (
	"encoding/json"

	dockerTypes "github.com/docker/docker/api/types"
	filters "github.com/docker/docker/api/types/filters"
	"github.com/google/go-querystring/query"
)

// GetNetwork get a network
func (api *API) GetNetwork(name string) ([]dockerTypes.NetworkResource, error) {
	path := "/networks"
	var networks []dockerTypes.NetworkResource
	arg := filters.NewArgs(filters.KeyValuePair{
		Key:   "name",
		Value: name,
	})
	// opt := dockerTypes.NetworkListOptions{Filters: arg}

	q, err := json.Marshal(arg)
	if err != nil {
		return networks, err
	}
	type Filters struct {
		Items string `url:"filters"`
	}
	filters := Filters{Items: string(q)}
	qs, err := query.Values(filters)
	if err != nil {
		return networks, err
	}

	if err := api.get(path, qs.Encode(), &networks); err != nil {
		return networks, err
	}

	return networks, nil
}

// CreateNetwork create a network
func (api *API) CreateNetwork(name string) error {
	path := "/networks/create"
	ncReq := dockerTypes.NetworkCreateRequest{
		Name: name,
	}

	body, err := json.Marshal(ncReq)
	if err != nil {
		return err
	}

	var resp dockerTypes.NetworkCreateResponse
	if err := api.post(path, body, 201, resp); err != nil {
		return err
	}
	return nil
}
