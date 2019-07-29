package commands

import (
	"log"

	"github.com/metrue/fx/config"
	"github.com/metrue/fx/utils"
)

// AddHost add a host
func (cmder *Commander) AddHost(name string, host config.Host) error {
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
	return cmder.cfg.AddHost(name, host)
}

// RemoveHost remote a host
func (cmder *Commander) RemoveHost(name string) error {
	return cmder.cfg.RemoveHost(name)
}

// ListHosts list hosts
func (cmder *Commander) ListHosts() error {
	hosts, err := cmder.cfg.ListHosts()
	if err != nil {
		return err
	}

	return utils.OutputJSON(hosts)
}

// SetDefaultHost set default host
func (cmder *Commander) SetDefaultHost(name string) error {
	host, err := cmder.cfg.GetHost(name)
	if err != nil {
		return err
	}
	return cmder.cfg.SetDefaultHost(name, host)
}

// GetDefaultHost get default host
func (cmder *Commander) GetDefaultHost() error {
	host, err := cmder.cfg.GetDefaultHost()
	if err != nil {
		return err
	}
	return utils.OutputJSON(host)
}
