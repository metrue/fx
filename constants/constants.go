package constants

import (
	"os"
	"path"
)

// AgentPort fx server agent port
const AgentPort = "8866"

// BaseImages base images to build fx functions
var BaseImages = []string{
	"metrue/fx-java-base",
	"metrue/fx-julia-base",
	"metrue/fx-python-base",
	"metrue/fx-node-base",
	"metrue/fx-d-base",
	"metrue/fx-go-base",
}

// ConfigPath path to config
var ConfigPath = path.Join(os.Getenv("HOME"), ".fx")

// AgentContainerName fx agent name
const AgentContainerName = "fx-agent"

// CheckedSymbol check symbol ✓
const CheckedSymbol = "\u2713"
