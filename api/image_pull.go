package api

import (
	"bufio"
	"fmt"
	"net/http"
	"time"
)

func (api *API) pull(name string) error {
	path := fmt.Sprintf("/%s/images/create?fromImage=%s&tag=latest", api.version, name)
	url := fmt.Sprintf("http://%s%s", api.endpoint, path)
	req, err := http.NewRequest("POST", url, nil)
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)

	for scanner.Scan() {
		fmt.Printf("%s\n", scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
