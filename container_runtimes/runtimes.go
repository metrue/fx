package containerruntimes

import (
	"context"

	"github.com/metrue/fx/types"
)

// ContainerRuntime interface
type ContainerRuntime interface {
	BuildImage(ctx context.Context, workdir string, name string) error
	PushImage(ctx context.Context, name string) (string, error)
	InspectImage(ctx context.Context, name string, img interface{}) error
	TagImage(ctx context.Context, name string, tag string) error
	StartContainer(ctx context.Context, name string, image string, bindings []types.PortBinding) error
	StopContainer(ctx context.Context, name string) error
	RemoveContainer(ctx context.Context, name string) error
	InspectContainer(ctx context.Context, name string, container interface{}) error
	ListContainer(ctx context.Context, filter string) ([]types.Service, error)
	Version(ctx context.Context) (string, error)
}
