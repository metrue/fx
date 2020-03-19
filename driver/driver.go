package driver

import (
	"context"

	"github.com/metrue/fx/types"
)

// Driver fx function running driver
type Driver interface {
	Deploy(ctx context.Context, fn string, name string, image string, bindings []types.PortBinding) error
	Destroy(ctx context.Context, name string) error
	Update(ctx context.Context, name string) error
	GetStatus(ctx context.Context, name string) (types.Service, error)
	List(ctx context.Context, name string) ([]types.Service, error)
	Ping(ctx context.Context) error
}
