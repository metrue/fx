package deploy

import "context"

// Deployer make a image a service
type Deployer interface {
	Deploy(ctx context.Context, workdir string, name string, ports []int32) error
	Destroy(ctx context.Context, name string) error
	Update(ctx context.Context, name string) error
	GetStatus(ctx context.Context, name string) error
}
