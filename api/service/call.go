package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/metrue/fx/api"
	"github.com/metrue/fx/pkg/utils"
)

func Call(ctx context.Context, req *api.CallRequest) (*api.CallResponse, error) {
	funcMeta := api.FunctionMeta{
		Lang:    req.Lang,
		Path:    req.Path,
		Content: req.Content,
	}

	upRes, err := DoUp(funcMeta)
	if err != nil {
		return nil, err
	}

	time.Sleep(time.Second * 2)

	params := utils.PairsToParams(strings.Fields(req.Params))
	body, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("http://%s", upRes.LocalAddress)
	r, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &api.CallResponse{
		Error: "",
		Data:  string(buf),
	}, nil
}
