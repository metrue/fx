package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path"

	dockerInfra "github.com/metrue/fx/infra/docker"
	"github.com/metrue/fx/types"
	"github.com/metrue/fx/utils"
	"github.com/mitchellh/go-homedir"
)

// Configer manage fx config
type Configer interface {
	GetCurrentCloud() ([]byte, error)
	GetCurrentCloudType() (string, error)
	GetKubeConfig() (string, error)
	UseCloud(name string) error
	View() ([]byte, error)
	AddCloud(name string, meta []byte) error
}

// Config config of fx
type Config struct {
	configFile string
	container  *Container
}

const defaultFxConfig = "~/.fx/config.yml"

// LoadDefault load default config
func LoadDefault() (*Config, error) {
	configFile, err := homedir.Expand(defaultFxConfig)
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
	}
	return load(configFile)
}

func load(configFile string) (*Config, error) {
	container, err := CreateContainer(configFile)
	if err != nil {
		return nil, err
	}
	config := &Config{
		configFile: configFile,
		container:  container,
	}

	if container.get("clouds") == nil {
		if err := config.writeDefaultConfig(); err != nil {
			return nil, err
		}
	}
	return config, nil
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
	}
	return load(configFile)
}

// AddCloud add k8s cloud
func (c *Config) AddCloud(name string, meta []byte) error {
	var cloudMeta map[string]interface{}
	if err := json.Unmarshal(meta, &cloudMeta); err != nil {
		return err
	}

	cloudType, ok := cloudMeta["type"].(string)
	if !ok || cloudType == "" {
		return fmt.Errorf("unknown cloud type")
	}

	if cloudType == types.CloudTypeK8S {
		dir := path.Dir(c.configFile)
		kubecfg := path.Join(dir, name+".kubeconfig")
		if err := utils.EnsureFile(kubecfg); err != nil {
			return err
		}
		config, ok := cloudMeta["config"].(string)
		if !ok {
			return fmt.Errorf("invalid k8s config")
		}
		if err := ioutil.WriteFile(kubecfg, []byte(config), 0666); err != nil {
			return err
		}
	}

	if err := c.container.set("clouds."+name, cloudMeta); err != nil {
		return err
	}

	return nil
}

// UseCloud set cloud instance with name as current context
func (c *Config) UseCloud(name string) error {
	if name == "" {
		return fmt.Errorf("could not use empty name")
	}

	if c.container.get("clouds."+name) == nil {
		return fmt.Errorf("no such cloud with name: %s", name)
	}
	return c.container.set("current_cloud", name)
}

// View view current config
func (c *Config) View() ([]byte, error) {
	return ioutil.ReadFile(c.configFile)
}

// GetCurrentCloud get current using cloud's meta
func (c *Config) GetCurrentCloud() ([]byte, error) {
	name, ok := c.container.get("current_cloud").(string)
	if !ok {
		return nil, fmt.Errorf("no active cloud")
	}
	meta := c.container.get("clouds." + name)
	if meta == nil {
		return nil, fmt.Errorf("invalid config")
	}
	return json.Marshal(meta)
}

// GetCurrentCloudType get current cloud type
func (c *Config) GetCurrentCloudType() (string, error) {
	name, ok := c.container.get("current_cloud").(string)
	if !ok {
		return "", fmt.Errorf("no active cloud")
	}
	return c.container.get("clouds." + name + ".type").(string), nil
}

// GetKubeConfig get kubeconfig
func (c *Config) GetKubeConfig() (string, error) {
	name, ok := c.container.get("current_cloud").(string)
	if !ok {
		return "", fmt.Errorf("no active cloud")
	}
	dir := path.Dir(c.configFile)
	kubecfg := path.Join(dir, name+".kubeconfig")
	return kubecfg, nil
}

func (c *Config) writeDefaultConfig() error {
	me, err := user.Current()
	if err != nil {
		return err
	}

	defaultCloud := &dockerInfra.Cloud{
		IP:   "127.0.0.1",
		User: me.Username,
		Name: "default",
		Type: types.CloudTypeDocker,
	}
	meta, err := defaultCloud.Dump()
	if err != nil {
		return err
	}
	if err := c.container.set("clouds", map[string]interface{}{}); err != nil {
		return err
	}
	if err := c.AddCloud("default", meta); err != nil {
		return err
	}
	return c.UseCloud("default")
}

var (
	_ Configer = &Config{}
)
