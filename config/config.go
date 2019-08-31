package config

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

// Configer interface
type Configer interface {
	GetMachine(name string) (Host, error)
	AddMachine(name string, host Host) error
	RemoveHost(name string) error
	ListActiveMachines() (map[string]Host, error)
	ListMachines() (map[string]Host, error)
	EnableMachine(name string) error
	DisableMachine(name string) error
	UpdateProvisionedStatus(name string, ok bool) error
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

		localhost := Host{
			Host:        "localhost",
			Password:    "",
			User:        "",
			Enabled:     true,
			Provisioned: false,
		}
		viper.Set("hosts", map[string]Host{"localhost": localhost})
		return viper.WriteConfig()
	}

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Fatal error config file: %s", err)
	}
	return nil
}

// GetMachine get host by name
func (c *Config) GetMachine(name string) (Host, error) {
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

// ListActiveMachines list enabled machines
func (c *Config) ListActiveMachines() (map[string]Host, error) {
	hosts, err := c.ListMachines()
	if err != nil {
		return map[string]Host{}, err
	}
	lst := map[string]Host{}
	for name, h := range hosts {
		if h.Enabled {
			lst[name] = h
		}
	}
	return lst, nil
}

// AddMachine add host
func (c *Config) AddMachine(name string, host Host) error {
	if !viper.IsSet("hosts") {
		viper.Set("hosts", map[string]Host{})
	}

	hosts, err := c.ListMachines()
	if err != nil {
		return err
	}
	hosts[name] = host
	viper.Set("hosts", hosts)
	return viper.WriteConfig()
}

// RemoveHost remote a host
func (c *Config) RemoveHost(name string) error {
	hosts, err := c.ListMachines()
	if err != nil {
		return err
	}

	if len(hosts) == 1 {
		return fmt.Errorf("only one host left now, at least one host required by fx")
	}

	if _, ok := hosts[name]; ok {
		delete(hosts, name)

		viper.Set("hosts", hosts)
		return viper.WriteConfig()
	}
	return fmt.Errorf("no such host %s", name)
}

// ListMachines list hosts
func (c *Config) ListMachines() (map[string]Host, error) {
	var hosts map[string]Host
	if err := viper.UnmarshalKey("hosts", &hosts); err != nil {
		return nil, err
	}
	return hosts, nil
}

// EnableMachine enable a machine, after machine enabled, function will be deployed onto it when ever `fx up` invoked
func (c *Config) EnableMachine(name string) error {
	host, err := c.GetMachine(name)
	if err != nil {
		return err
	}
	host.Enabled = true

	if !viper.IsSet("hosts") {
		viper.Set("hosts", map[string]Host{})
	}

	hosts, err := c.ListMachines()
	if err != nil {
		return err
	}
	hosts[name] = host
	viper.Set("hosts", hosts)
	return viper.WriteConfig()
}

// DisableMachine disable a machine, after machine disabled, function will not be deployed onto it
func (c *Config) DisableMachine(name string) error {
	host, err := c.GetMachine(name)
	if err != nil {
		return err
	}
	host.Enabled = false

	if !viper.IsSet("hosts") {
		viper.Set("hosts", map[string]Host{})
	}

	hosts, err := c.ListMachines()
	if err != nil {
		return err
	}
	hosts[name] = host
	viper.Set("hosts", hosts)
	return viper.WriteConfig()
}

// UpdateProvisionedStatus update provisioned status
func (c *Config) UpdateProvisionedStatus(name string, ok bool) error {
	host, err := c.GetMachine(name)
	if err != nil {
		return err
	}
	host.Provisioned = ok

	if !viper.IsSet("hosts") {
		viper.Set("hosts", map[string]Host{})
	}

	hosts, err := c.ListMachines()
	if err != nil {
		return err
	}
	hosts[name] = host
	viper.Set("hosts", hosts)
	return viper.WriteConfig()
}

// IsMachineProvisioned check if machine provisioned
func (c *Config) IsMachineProvisioned(name string) bool {
	host, err := c.GetMachine(name)
	if err != nil {
		return false
	}
	return host.Provisioned
}
