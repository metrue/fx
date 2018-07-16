package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/metrue/fx/common"
)

type CallOutput struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func Call(address string, function string, params map[string]string) error {
	res, err := InvokeUpRequest(address, []string{function})
	if err != nil {
		common.HandleCallResult(CallOutput{Error: err.Error()})
		return nil
	}

	if len(res.Instances) != 1 {
		common.HandleCallResult(CallOutput{
			Error: fmt.Sprintf("could not up function: %s", function),
		})
		return nil
	}
	defer Down(address, []string{res.Instances[0].FunctionID})

	body, err := json.Marshal(params)
	if err != nil {
		common.HandleCallResult(CallOutput{Error: err.Error()})
		return nil
	}

	time.Sleep(time.Second * 2)

	url := fmt.Sprintf("http://%s", res.Instances[0].LocalAddress)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		common.HandleCallResult(CallOutput{Error: err.Error()})
		return nil
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		common.HandleCallResult(CallOutput{Error: err.Error()})
		return err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		common.HandleCallResult(CallOutput{Error: err.Error()})
	} else {
		common.HandleCallResult(CallOutput{Message: string(buf)})
	}

	return nil
}
