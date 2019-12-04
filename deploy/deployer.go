package deploy

import (
	"context"

	types "github.com/metrue/fx/types"
)

// Deployer make a image a service
type Deployer interface {
	Deploy(ctx context.Context, fn types.Func, name string, image string, bindings []types.PortBinding) error
	Destroy(ctx context.Context, name string) error
	Update(ctx context.Context, name string) error
	GetStatus(ctx context.Context, name string) (types.Service, error)
	List(ctx context.Context, name string) ([]types.Service, error)
	Ping(ctx context.Context) error
}
