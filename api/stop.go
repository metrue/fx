package api

import (
	"fmt"
	"net/http"
	"time"
)

// Stop a container by name
func (api *API) Stop(name string) error {
	path := fmt.Sprintf("/containers/%s/stop", name)
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

	return nil
}
