package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
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

// HasDockerfile check if there is Dockerfile in dir
func HasDockerfile(dir string) bool {
	var dockerfile string
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// nolint
		if info.Mode().IsRegular() && info.Name() == "Dockerfile" {
			dockerfile = path
		}
		return nil
	}); err != nil {
		return false
	}
	if dockerfile == "" {
		return false
	}
	return true
}
