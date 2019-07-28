package config

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/viper"
)

// Init config
func Init(configPath string) error {
	os.MkdirAll(configPath, os.ModePerm)

	ext := "yaml"
	name := "config"
	viper.SetConfigType(ext)
	viper.SetConfigName(name)
	viper.AddConfigPath(configPath)

	// detect if file exists
	configFilePath := path.Join(configPath, name+"."+ext)
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
func GetHost(name string) (Host, error) {
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
func GetDefaultHost() (Host, error) {
	var host Host
	if err := viper.UnmarshalKey("default", &host); err != nil {
		return Host{}, err
	}
	return host, nil
}

// SetDefaultHost set default host
// TODO no need name
func SetDefaultHost(name string, host Host) error {
	viper.Set("default", host)
	return viper.WriteConfig()
}

// IsRemote if running on remote
func IsRemote() bool {
	host, err := GetDefaultHost()
	if err != nil {
		return false
	}
	return host.IsRemote()
}

// AddHost add host
func AddHost(name string, host Host) error {
	if !viper.IsSet("hosts") {
		viper.Set("hosts", map[string]Host{})
	}

	hosts, err := ListHosts()
	if err != nil {
		return err
	}
	hosts[name] = host
	viper.Set("hosts", hosts)
	return viper.WriteConfig()
}

// RemoveHost remote a host
func RemoveHost(name string) error {
	hosts, err := ListHosts()
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
func ListHosts() (map[string]Host, error) {
	var hosts map[string]Host
	if err := viper.UnmarshalKey("hosts", &hosts); err != nil {
		return nil, err
	}
	return hosts, nil
}
