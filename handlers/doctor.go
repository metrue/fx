package handlers

import (
	"os"

	"github.com/apex/log"
	"github.com/metrue/fx/constants"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/doctor"
)

// Doctor command handle
func Doctor(ctx *context.Context) error {
	host := os.Getenv("DOCKER_REMOTE_HOST_ADDR")
	user := os.Getenv("DOCKER_REMOTE_HOST_USER")
	password := os.Getenv("DOCKER_REMOTE_HOST_PASSWORD")
	if host == "" {
		host = "localhost"
	}
	if err := doctor.New(host, user, password).Start(); err != nil {
		log.Warnf("machine %s is in dirty state: %v", host, err)
	} else {
		log.Infof("machine %s is in healthy state: %s", host, constants.CheckedSymbol)
	}
	return nil
}
