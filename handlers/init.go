package handlers

import (
	"os"

	"github.com/apex/log"
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/provision"
)

// Init start fx-agent
func Init() HandleFunc {
	return func(ctx *context.Context) error {
		host := os.Getenv("DOCKER_REMOTE_HOST_ADDR")
		user := os.Getenv("DOCKER_REMOTE_HOST_USER")
		passord := os.Getenv("DOCKER_REMOTE_HOST_PASSWORD")
		if host == "" {
			host = "127.0.0.1"
		}
		provisioner := provision.NewWithHost(host, user, passord)
		if !provisioner.IsFxAgentRunning() {
			if err := provisioner.StartFxAgent(); err != nil {
				log.Fatalf("could not start fx agent on host: %s", err)
				return err
			}
			log.Info("fx agent started")
		}
		log.Info("fx agent already started")
		return nil
	}
}
