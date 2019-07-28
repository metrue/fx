package commands

import (
	"log"

	"github.com/metrue/fx/config"
	"github.com/metrue/fx/utils"
)

// AddHost add a host
func AddHost(name string, host config.Host) error {
	if !host.Valid() {
		log.Fatalf("invaid host %v", host)
		return nil
	}

	if host.IsRemote() {
		if host.User == "" || host.Password == "" {
			log.Fatalf("the host to add is a remote, user and password for SSH login is required")
			return nil
		}
	}
	return config.AddHost(name, host)
}

// RemoveHost remote a host
func RemoveHost(name string) error {
	return config.RemoveHost(name)
}

// ListHosts list hosts
func ListHosts() error {
	hosts, err := config.ListHosts()
	if err != nil {
		return err
	}

	return utils.OutputJSON(hosts)
}

// SetDefaultHost set default host
func SetDefaultHost(name string) error {
	host, err := config.GetHost(name)
	if err != nil {
		return err
	}
	return config.SetDefaultHost(name, host)
}

// GetDefaultHost get default host
func GetDefaultHost() error {
	host, err := config.GetDefaultHost()
	if err != nil {
		return err
	}
	return utils.OutputJSON(host)
}
