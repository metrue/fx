package config

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

// Configer interface
type Configer interface {
	GetDefaultHost() (Host, error)
}

// Config config of fx
type Config struct {
	dir string
}

// New create a config
func New(dir string) *Config {
	return &Config{dir: dir}
}

// Init config
func (c *Config) Init() error {
	if err := os.MkdirAll(c.dir, os.ModePerm); err != nil {
		return err
	}

	ext := "yaml"
	name := "config"
	viper.SetConfigType(ext)
	viper.SetConfigName(name)
	viper.AddConfigPath(c.dir)

	// detect if file exists
	configFilePath := path.Join(c.dir, name+"."+ext)
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		fd, err := os.Create(configFilePath)
		if err != nil {
			return err
		}
		fd.Close()

		viper.Set("default", Host{
			Host:     "localhost",
			Password: "",
			User:     "",
		})

		viper.Set("hosts", map[string]Host{
			"localhost": Host{
				Host:     "localhost",
				Password: "",
				User:     "",
			},
		})
		return viper.WriteConfig()
	}

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Fatal error config file: %s", err)
	}
	return nil
}

// GetHost get host by name
func (c *Config) GetHost(name string) (Host, error) {
	var hosts map[string]Host
	if err := viper.UnmarshalKey("hosts", &hosts); err != nil {
		return Host{}, err
	}
	host, ok := hosts[name]
	if !ok {
		return Host{}, fmt.Errorf("no such host %v", name)
	}
	return host, nil
}

// GetDefaultHost get host
func (c *Config) GetDefaultHost() (Host, error) {
	var host Host
	if err := viper.UnmarshalKey("default", &host); err != nil {
		return Host{}, err
	}
	return host, nil
}

// SetDefaultHost set default host
// TODO no need name
func (c *Config) SetDefaultHost(name string, host Host) error {
	viper.Set("default", host)
	return viper.WriteConfig()
}

// IsRemote if running on remote
func (c *Config) IsRemote() bool {
	host, err := c.GetDefaultHost()
	if err != nil {
		return false
	}
	return host.IsRemote()
}

// AddHost add host
func (c *Config) AddHost(name string, host Host) error {
	if !viper.IsSet("hosts") {
		viper.Set("hosts", map[string]Host{})
	}

	hosts, err := c.ListHosts()
	if err != nil {
		return err
	}
	hosts[name] = host
	viper.Set("hosts", hosts)
	return viper.WriteConfig()
}

// RemoveHost remote a host
func (c *Config) RemoveHost(name string) error {
	hosts, err := c.ListHosts()
	if err != nil {
		return err
	}

	if len(hosts) == 1 {
		return fmt.Errorf("only one host left now, at least one host required by fx")
	}

	if _, ok := hosts[name]; ok {
		delete(hosts, name)
		return nil
	}
	return fmt.Errorf("no such host %s", name)
}

// ListHosts list hosts
func (c *Config) ListHosts() (map[string]Host, error) {
	var hosts map[string]Host
	if err := viper.UnmarshalKey("hosts", &hosts); err != nil {
		return nil, err
	}
	return hosts, nil
}
