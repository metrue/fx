package commands

import "github.com/metrue/fx/config"

func Use(address string) error {
	config := config.GetConfig()
	config.HttpServerAddr = address
	config.GrpcEndpoint = address
	return config.Save()
}
