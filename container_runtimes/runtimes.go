package containerruntimes

import (
	"context"

	"github.com/metrue/fx/types"
)

// ContainerRuntime interface
type ContainerRuntime interface {
	BuildImage(ctx context.Context, workdir string, name string) error
	PushImage(ctx context.Context, name string) (string, error)
	InspectImage(ct context.Context, name string, img interface{}) error
	StartContainer(ctx context.Context, name string, image string, bindings []types.PortBinding) error
	StopContainer(ctx context.Context, name string) error
	InspectContainer(ctx context.Context, name string, container interface{}) error
}
