package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func Call(address string, function string, params map[string]interface{}) error {
	id, service, err := Up(address, []string{function})
	if err != nil {
		return err
	}
	defer Down(address, []string{id})

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	time.Sleep(time.Second * 2)

	url := fmt.Sprintf("http://%s", service)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(buf))
	return nil
}
