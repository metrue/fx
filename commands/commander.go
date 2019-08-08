package commands

import "github.com/metrue/fx/config"

// Commandor interface
type Commandor interface {
}

// Commander commands management
type Commander struct {
	cfg *config.Config
}

// New create commander
func New(cfg *config.Config) *Commander {
	return &Commander{
		cfg: cfg,
	}
}
