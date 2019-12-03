package middlewares

import (
	"github.com/metrue/fx/config"
	"github.com/metrue/fx/context"
)

// LoadConfig load default config
func LoadConfig(ctx context.Contexter) error {
	config, err := config.LoadDefault()
	if err != nil {
		return err
	}
	ctx.Set("config", config)
	return nil
}
