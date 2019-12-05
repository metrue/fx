package infra

import (
	"context"

	"github.com/metrue/fx/types"
)

// Provisioner provision interface
type Provisioner interface {
	Provision() (config []byte, err error)
	HealthCheck() (bool, error)
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
	Provisioner
	Deployer
}
