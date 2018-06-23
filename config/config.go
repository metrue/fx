package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"

	"github.com/metrue/fx/pkg/utils"
)

type FxConfig struct {
	HttpServerAddr string `json:"http"`
	GrpcEndpoint   string `json:"grpc"`
}

var CONFIG = path.Join(os.Getenv("HOME"), ".fx/config.json")

func GetConfig() *FxConfig {
	defaultConfig := &FxConfig{
		HttpServerAddr: "localhost:30080",
		GrpcEndpoint:   "localhost:5000",
	}

	raw, err := ioutil.ReadFile(CONFIG)
	if err != nil {
		return defaultConfig
	}

	var c FxConfig
	json.Unmarshal(raw, &c)

	if len(c.HttpServerAddr) > 0 && len(c.GrpcEndpoint) > 0 {
		return &c
	}

	return defaultConfig
}

func (c *FxConfig) Save() error {
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

func GetGrpcEndpoint() string {
	c := GetConfig()
	return c.HttpServerAddr
}

func GetHttpServerAddr() string {
	c := GetConfig()
	return c.GrpcEndpoint
}
