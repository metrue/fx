package provisioners

import "context"

// Provisioner define provisioner interface
type Provisioner interface {
	Provision(ctx context.Context, isRemote bool) error
}
