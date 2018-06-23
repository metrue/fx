package commands

import "github.com/metrue/fx/config"

func Use(address string) error {
	config := config.GetConfig()
	return config.SetHost(address)
}
