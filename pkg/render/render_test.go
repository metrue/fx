package render

import (
	"testing"

	"github.com/metrue/fx/types"
)

func TestTable(t *testing.T) {
	services := []types.Service{
		types.Service{
			ID:   "id-1",
			Name: "name-1",
			Host: "127.0.0.1",
			Port: 1000,
		},
	}
	Table(services)
}
