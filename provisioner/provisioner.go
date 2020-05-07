package provisioner

import (
	"context"
	"time"
)

// SSHConnectionTimeout default timeout for ssh connection
const SSHConnectionTimeout = 10 * time.Second

// Provisioner define provisioner interface
type Provisioner interface {
	Provision(ctx context.Context, isRemote bool) error
}
