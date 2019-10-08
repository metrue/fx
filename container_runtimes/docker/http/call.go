package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/apex/log"
	"github.com/metrue/fx/types"
	"github.com/metrue/fx/utils"
)

// Call function directly with given params
func (api *API) Call(file string, param string, project types.Project) error {
	service, err := api.Build(project)
	if err != nil {
		log.Fatalf("Build Service: %v", err)
		return err
	}
	log.Info("Build Service: \u2713")

	if err := api.Run(9999, &service); err != nil {
		log.Fatalf("Run Service: %v", err)
		return err
	}
	log.Info("Run Service: \u2713")

	params := utils.PairsToParams(strings.Fields(param))
	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	// Wait 2 seconds for service startup
	time.Sleep(time.Second * 2)

	url := fmt.Sprintf("http://%s:%d", service.Host, service.Port)
	r, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	r.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(r)
	if err != nil {
		log.Fatalf("Call Service: %v", err)
		return err
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Call Service: %v", err)
		return err
	}
	log.Info("Call Service: \u2713")
	return utils.OutputJSON(string(buf))
}
