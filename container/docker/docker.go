package docker

import (
	"context"

	"github.com/metrue/fx/container"
)

type Docker struct {
}

func (d *Docker) Deploy(ctx context.Context, name string, image string, port []int32) error {
	return nil
}

func (d *Docker) Update(ctx context.Context, name string) error {
	return nil
}

func (d *Docker) Destroy(ctx context.Context, name string) error {
	return nil
}

func (d *Docker) GetStatus(ctx context.Context, name string) error {
	return nil
}

var (
	_ container.Runner = &Docker{}
)
