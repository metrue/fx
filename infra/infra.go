package infra

import (
	"context"

	"github.com/metrue/fx/types"
)

// Clouder cloud interface
type Clouder interface {
	Provision() error
	GetConfig() (string, error)
	GetType() string
	Dump() ([]byte, error)
}

// Deployer deploy interface
type Deployer interface {
	Deploy(ctx context.Context, fn string, name string, image string, bindings []types.PortBinding) error
	Destroy(ctx context.Context, name string) error
	Update(ctx context.Context, name string) error
	GetStatus(ctx context.Context, name string) (types.Service, error)
	List(ctx context.Context, name string) ([]types.Service, error)
	Ping(ctx context.Context) error
}

// Infra infrastructure provision interface
type Infra interface {
	Deployer
}
