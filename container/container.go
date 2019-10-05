package container

import "context"

// Runner make a image a service
type Runner interface {
	Deploy(ctx context.Context, name string, image string, ports []int32) error
	Destroy(ctx context.Context, name string) error
	Update(ctx context.Context, name string) error
	GetStatus(ctx context.Context, name string) error
}
