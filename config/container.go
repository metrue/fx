package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/metrue/fx/utils"
	"github.com/spf13/viper"
)

// Container config container, wrap viper as a key-value store with lock
type Container struct {
	mux   sync.Mutex
	store string
}

// CreateContainer new a container
func CreateContainer(storeFile string) (*Container, error) {
	if err := utils.EnsureFile(storeFile); err != nil {
		return nil, err
	}

	dir := filepath.Dir(storeFile)
	ext := filepath.Ext(storeFile)
	name := filepath.Base(storeFile)
	viper.AddConfigPath(dir)
	viper.SetConfigName(strings.Replace(name, ext, "", 1))
	viper.SetConfigType(strings.Replace(ext, ".", "", 1))
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	return &Container{
		store: storeFile,
	}, nil

}

func (c *Container) set(key string, value interface{}) error {
	c.mux.Lock()
	defer c.mux.Unlock()

	if key == "" {
		return fmt.Errorf("empty key not allowed")
	}

	keys := strings.Split(key, ".")
	if len(keys) == 1 {
		viper.Set(key, value)
	} else {
		prePath := keys[0]
		for i := 1; i < len(keys)-2; i++ {
			prePath += "." + keys[i]
		}
		if viper.Get(prePath) == nil {
			return fmt.Errorf("%s not existed", prePath)
		}
		viper.Set(key, value)
	}
	// viper.Set(key, value)
	if err := viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}

func (c *Container) get(key string) interface{} {
	c.mux.Lock()
	defer c.mux.Unlock()

	return viper.Get(key)
}
