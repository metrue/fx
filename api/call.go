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
	"github.com/google/uuid"
	"github.com/metrue/fx/types"
	"github.com/metrue/fx/utils"
)

// Call function directly with given params
func (api *API) Call(file string, param string) error {
	src, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatalf("Read Source: %v", err)
		return err
	}
	log.Info("Read Source: \u2713")

	lang := utils.GetLangFromFileName(file)
	fn := types.ServiceFunctionSource{
		Language: lang,
		Source:   string(src),
	}

	project, err := api.Pack(uuid.New().String(), fn)
	if err != nil {
		log.Fatalf("Pack Service: %v", err)
		return err
	}
	log.Info("Pack Service: \u2713")

	service, err := api.Build(project)
	if err != nil {
		log.Fatalf("Build Service: %v", err)
		return err
	}
	log.Info("Build Service: \u2713")

	if err := api.Run(&service); err != nil {
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

	url := fmt.Sprintf("http://%s:%d", service.Instances[0].Host, service.Instances[0].Port)
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
	utils.OutputJSON(string(buf))

	return nil
}
