package middlewares

import (
	"github.com/metrue/fx/context"
	"github.com/metrue/fx/hook"
)

// Hook midlleware
func Hook() func(ctx context.Contexter) (err error) {
	return func(ctx context.Contexter) error {
		hooks, err := hook.Descovery("")
		if err != nil {
			return err
		}
		ctx.Set("hooks", hooks)
		return nil
	}
}
