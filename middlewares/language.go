package middlewares

import (
	"fmt"
	"path/filepath"

	"github.com/metrue/fx/constants"
	"github.com/metrue/fx/context"
)

// Language to find out what language of function is
func Language() func(ctx context.Contexter) (err error) {
	return func(ctx context.Contexter) error {
		fn := ctx.Get("fn").(string)
		ext := filepath.Ext(fn)
		language, ok := constants.ExtLangMapping[ext]
		if !ok {
			return fmt.Errorf("%s not supported yet", ext)
		}
		ctx.Set("language", language)
		return nil
	}
}
