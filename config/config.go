package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/metrue/fx/utils"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// Cloud define cloud infrastructure
type Cloud map[string]string

// Config config of fx
type Config struct {
	Clouds       map[string]Cloud
	CurrentCloud string
}

// New create a config
func New() *Config {
	return &Config{}
}

// Load config
func Load() (*Config, error) {
	configFile, err := getConfigFile()
	if err != nil {
		return nil, err
	}

	ext := "yaml"
	name := "config"
	viper.SetConfigType(ext)
	viper.SetConfigName(name)
	viper.SetConfigFile(configFile)

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := utils.EnsureFile(configFile); err != nil {
			return nil, err
		}
		writeDefaultConfig()
	}

	return load()
}

// AddCloud add a cloud
func AddCloud(name string, cloud Cloud) error {
	config, err := load()
	if err != nil {
		return err
	}

	config.Clouds[name] = cloud

	return save(config)
}

// Use set cloud instance with name as current context
func Use(name string) error {
	config, err := load()
	if err != nil {
		return err
	}

	has := false
	for n := range config.Clouds {
		if n == name {
			has = true
			break
		}
	}
	if !has {
		return fmt.Errorf("no cloud with name = %s", name)
	}
	config.CurrentCloud = name
	return nil
}

// View view current config
func View() ([]byte, error) {
	configFile, err := getConfigFile()
	if err != nil {
		return []byte{}, err
	}
	return ioutil.ReadFile(configFile)
}

func load() (*Config, error) {
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	config := &Config{}

	config.CurrentCloud = viper.GetString("current_cloud")
	clouds := make(map[string]Cloud)

	cloudList := viper.Get("clouds").(map[string]Cloud)
	for n, c := range cloudList {
		clouds[n] = c
	}
	config.Clouds = clouds

	return config, nil
}

func save(c *Config) error {
	viper.Set("clouds", c.Clouds)
	viper.Set("current_cloud", c.CurrentCloud)
	return viper.WriteConfig()
}

func getConfigFile() (string, error) {
	configFile, err := homedir.Expand("~/.fx/config.yml")
	if err != nil {
		return "", err
	}

	if os.Getenv("FX_CONFIG") != "" {
		configFile = os.Getenv("FX_CONFIG")
	}
	return configFile, nil
}

func writeDefaultConfig() error {
	me, err := user.Current()
	if err != nil {
		return err
	}
	viper.Set("current_cloud", "default")
	viper.Set("clouds", map[string]Cloud{
		"default": Cloud{
			"type": "docker",
			"host": "127.0.0.1",
			"user": me.Username,
		},
	})
	if err := viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}
