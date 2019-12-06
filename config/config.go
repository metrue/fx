package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"
	"sync"

	"github.com/metrue/fx/utils"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
)

// Items data of config file
type Items struct {
	Clouds       map[string]map[string]string `json:"clouds"`
	CurrentCloud string                       `json:"current_cloud"`
}

// Config config of fx
type Config struct {
	mux        sync.Mutex
	configFile string
	Items
}

// LoadDefault load default config
func LoadDefault() (*Config, error) {
	configFile, err := homedir.Expand("~/.fx/config.yml")
	if err != nil {
		return nil, err
	}
	if os.Getenv("FX_CONFIG") != "" {
		configFile = os.Getenv("FX_CONFIG")

	}
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := utils.EnsureFile(configFile); err != nil {
			return nil, err
		}
		if err := writeDefaultConfig(configFile); err != nil {
			return nil, err
		}
	}
	return load(configFile)
}

// Load config
func Load(configFile string) (*Config, error) {
	if configFile == "" {
		return nil, fmt.Errorf("invalid config file")
	}

	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := utils.EnsureFile(configFile); err != nil {
			return nil, err
		}
		if err := writeDefaultConfig(configFile); err != nil {
			return nil, err
		}
	}
	return load(configFile)
}

// AddCloud add a cloud
func (c *Config) addCloud(name string, cloud map[string]string) error {
	c.Items.Clouds[name] = cloud
	return save(c)
}

// AddDockerCloud add docker cloud
func (c *Config) AddDockerCloud(name string, config []byte) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	var conf map[string]string
	err := json.Unmarshal(config, &conf)
	if err != nil {
		return err
	}

	cloud := map[string]string{
		"type": "docker",
		"host": conf["ip"],
		"user": conf["user"],
	}
	return c.addCloud(name, cloud)
}

// AddK8SCloud add k8s cloud
func (c *Config) AddK8SCloud(name string, kubeconfig []byte) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	dir := path.Dir(c.configFile)
	kubecfg := path.Join(dir, name+".kubeconfig")
	if err := utils.EnsureFile(kubecfg); err != nil {
		return err
	}
	if err := ioutil.WriteFile(kubecfg, kubeconfig, 0666); err != nil {
		return err
	}

	cloud := map[string]string{
		"type":       "k8s",
		"kubeconfig": kubecfg,
	}

	return c.addCloud(name, cloud)
}

// Use set cloud instance with name as current context
func (c *Config) Use(name string) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	has := false
	for n := range c.Clouds {
		if n == name {
			has = true
			break
		}
	}
	if !has {
		return fmt.Errorf("no cloud with name = %s", name)
	}
	c.Items.CurrentCloud = name
	return save(c)
}

// View view current config
func (c *Config) View() ([]byte, error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	return ioutil.ReadFile(c.configFile)
}

func load(configFile string) (*Config, error) {
	conf, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var items Items
	if err := yaml.Unmarshal(conf, &items); err != nil {
		return nil, err
	}
	var c = Config{
		configFile: configFile,
		Items:      items,
	}
	return &c, nil
}

func save(c *Config) error {
	conf, err := yaml.Marshal(c.Items)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(c.configFile, conf, 0666); err != nil {
		return err
	}
	return nil
}

func writeDefaultConfig(configFile string) error {
	me, err := user.Current()
	if err != nil {
		return err
	}
	items := Items{
		Clouds: map[string]map[string]string{
			"default": map[string]string{
				"type": "docker",
				"host": "127.0.0.1",
				"user": me.Username,
			},
		},
		CurrentCloud: "default",
	}

	body, err := yaml.Marshal(items)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(configFile, body, 0666); err != nil {
		return err
	}

	return nil
}
