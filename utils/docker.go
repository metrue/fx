package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	dockerTypes "github.com/docker/docker/api/types"
)

// DockerVersion docker verion
func DockerVersion(host string, port string) (string, error) {
	path := host + ":" + port + "/version"
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
