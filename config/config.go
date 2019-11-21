package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/metrue/fx/utils"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

// Config config of fx
type Config struct {
	Clouds       map[string]interface{} `json:"clouds"`
	CurrentCloud string                 `json:"current_cloud"`
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

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := utils.EnsureFile(configFile); err != nil {
			return nil, err
		}
		if err := writeDefaultConfig(configFile); err != nil {
			return nil, err
		}
	}

	return load()
}

// AddCloud add a cloud
func AddCloud(name string, cloud interface{}) error {
	config, err := load()
	if err != nil {
		return err
	}

	config.Clouds[name] = cloud

	return save(config)
}

// AddDockerCloud add docker cloud
func AddDockerCloud(name string, host string, user string) error {
	cloud := DockerCloud{
		Host: host,
		User: user,
	}
	return AddCloud(name, cloud)
}

// AddK8SCloud add k8s cloud
func AddK8SCloud(name string, kubeconfig []byte) error {
	configFile, err := homedir.Expand("~/.fx/" + name + ".kubeconfig")
	if err != nil {
		return err
	}
	if err := utils.EnsureFile(configFile); err != nil {
		return err
	}
	if err := ioutil.WriteFile(configFile, kubeconfig, 0666); err != nil {
		return err
	}

	cloud := K8SCloud{
		KubeConfig: configFile,
	}

	return AddCloud(name, cloud)
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
	configFile, err := getConfigFile()
	if err != nil {
		return nil, err
	}
	conf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var c Config
	if err := yaml.Unmarshal(conf, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

func save(c *Config) error {
	conf, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	configFile, err := getConfigFile()
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(configFile, conf, 0666); err != nil {
		return err
	}
	return nil
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

func writeDefaultConfig(configFile string) error {
	me, err := user.Current()
	if err != nil {
		return err
	}
	conf := Config{
		Clouds: map[string]interface{}{
			"default": DockerCloud{
				Host: "127.0.0.1",
				User: me.Username,
			},
		},
		CurrentCloud: "default",
	}

	y, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(configFile, y, 0666); err != nil {
		return err
	}

	return nil
}
