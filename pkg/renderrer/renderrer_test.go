package renderrer

import (
	"testing"

	"github.com/metrue/fx/types"
)

func TestRenderrer(t *testing.T) {
	services := []types.Service{
		types.Service{
			ID:   "id-1",
			Name: "name-1",
			Host: "127.0.0.1",
			Port: 1000,
		},
	}
	t.Run("toTable", func(t *testing.T) {
		Render(services, "table")
	})
	t.Run("toJSON", func(t *testing.T) {
		Render(services, "json")
	})
}
