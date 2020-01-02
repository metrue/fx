package config

import (
	"os"
)

// DisableContainerAutoremove to tell if to run container with --rm
var DisableContainerAutoremove = false

func init() {
	if os.Getenv("DISABLE_CONTAINER_AUTOREMOVE") == "true" {
		DisableContainerAutoremove = true
	}
}
