package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/metrue/fx/pkg/utils"
)

type FxConfig struct {
	Http string `json:"http"`
	Grpc string `json:"grpc"`
}

var CONFIG = path.Join(os.Getenv("HOME"), ".fx/config.json")
var HTTP_PORT = 30080
var GRPC_PORT = 50000
var DEFAULT_CONFIG = &FxConfig{
	Http: fmt.Sprintf("localhost:%d", HTTP_PORT),
	Grpc: fmt.Sprintf("localhost:%d", GRPC_PORT),
}

func GetConfig() *FxConfig {
	_, err := os.Stat(CONFIG)
	if err != nil {
		DEFAULT_CONFIG.save()
		return DEFAULT_CONFIG
	}

	raw, err := ioutil.ReadFile(CONFIG)
	if err != nil {
		panic(err)
	}

	var c FxConfig
	json.Unmarshal(raw, &c)

	if len(c.Http) > 0 && len(c.Grpc) > 0 {
		return &c
	}

	return DEFAULT_CONFIG
}

func (c *FxConfig) save() error {
	os.Remove(CONFIG)

	utils.EnsureFile(CONFIG)

	configFile, err := os.OpenFile(CONFIG, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	configContent, err := json.Marshal(c)
	if err != nil {
		return err
	}

	_, err = configFile.Write(configContent)
	if err != nil {
		return err
	}

	return nil
}

func (c *FxConfig) SetHost(addr string) error {
	c.Http = fmt.Sprintf("%s:%d", addr, HTTP_PORT)
	c.Grpc = fmt.Sprintf("%s:%d", addr, GRPC_PORT)

	err := c.save()
	if err != nil {
		panic(err)
	}
	return err
}

func GetGrpcEndpoint() string {
	c := GetConfig()
	return c.Grpc
}

func GetHttpServerAddr() string {
	c := GetConfig()
	return c.Http
}
