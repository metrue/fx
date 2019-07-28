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

		viper.Set("host", "localhost")
		return viper.WriteConfig()
	}

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("Fatal error config file: %s", err)
	}
	return nil
}

// SetHost set host
func SetHost(host string) error {
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("config file not found: %s", err)
		}
		return err
	}

	viper.Set("host", host)
	return viper.WriteConfig()
}

// GetHost get host
func GetHost() string {
	return viper.GetString("host")
}

// IsRemote if running on remote
func IsRemote() bool {
	host := GetHost()
	return host != "127.0.0.1" && host != "localhost"
}
