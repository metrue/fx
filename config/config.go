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
	httpServerAddr string `json:"http"`
	grpcEndpoint   string `json:"grpc"`
}

var CONFIG = path.Join(os.Getenv("HOME"), ".fx/config.json")
var HTTP_PORT = 30080
var GRPC_PORT = 50000
var DEFAULT_CONFIG = &FxConfig{
	httpServerAddr: fmt.Sprintf("localhost:%d", HTTP_PORT),
	grpcEndpoint:   fmt.Sprintf("localhost:%d", GRPC_PORT),
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

	if len(c.httpServerAddr) > 0 && len(c.grpcEndpoint) > 0 {
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
	c.httpServerAddr = fmt.Sprintf("%s:30080", addr)
	c.grpcEndpoint = fmt.Sprintf("%s:5000", addr)

	return c.save()
}

func GetGrpcEndpoint() string {
	c := GetConfig()
	return c.grpcEndpoint
}

func GetHttpServerAddr() string {
	c := GetConfig()
	return c.httpServerAddr
}
