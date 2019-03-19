package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	dockerTypes "github.com/docker/docker/api/types"
)

// Version get version of dockerd server
func Version(endpoint string) (string, error) {
	path := "/version"
	url := fmt.Sprintf("%s%s", endpoint, path)
	if !strings.HasPrefix(url, "http") {
		url = "http://" + url
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("request %s failed: %d - %s", url, resp.StatusCode, resp.Status)
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
