package docker

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	dockerMock "github.com/metrue/fx/container_runtimes/mocks"
	"github.com/metrue/fx/types"
)

func TestDriverPing(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dockerClient := dockerMock.NewMockContainerRuntime(ctrl)
	n := New(Options{
		DockerClient: dockerClient,
	})
	ctx := context.Background()
	dockerClient.EXPECT().Version(ctx).Return("", nil)
	if err := n.Ping(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestDriverDeploy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dockerClient := dockerMock.NewMockContainerRuntime(ctrl)
	n := New(Options{
		DockerClient: dockerClient,
	})
	ctx := context.Background()
	fn := "fn"
	name := "name"
	image := "image"
	ports := []types.PortBinding{}
	dockerClient.EXPECT().StartContainer(ctx, name, image, ports).Return(nil)
	if err := n.Deploy(ctx, fn, name, image, ports); err != nil {
		t.Fatal(err)
	}
}

func TestDriverDestroy(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dockerClient := dockerMock.NewMockContainerRuntime(ctrl)
	n := New(Options{
		DockerClient: dockerClient,
	})
	ctx := context.Background()
	name := "name"
	dockerClient.EXPECT().StopContainer(ctx, name).Return(nil)
	dockerClient.EXPECT().RemoveContainer(ctx, name).Return(nil)
	if err := n.Destroy(ctx, name); err != nil {
		t.Fatal(err)
	}
}

func TestDriverGetStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dockerClient := dockerMock.NewMockContainerRuntime(ctrl)
	n := New(Options{
		DockerClient: dockerClient,
	})
	ctx := context.Background()
	name := "name"
	err := errors.New("no such container")
	dockerClient.EXPECT().InspectContainer(ctx, name, gomock.Any()).Return(err)
	if _, err := n.GetStatus(ctx, name); err == nil {
		t.Fatalf("should get error")
	}
}

func TestList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	dockerClient := dockerMock.NewMockContainerRuntime(ctrl)
	n := New(Options{
		DockerClient: dockerClient,
	})
	ctx := context.Background()
	name := "name"
	dockerClient.EXPECT().ListContainer(ctx, name).Return(nil, nil)
	if _, err := n.List(ctx, name); err != nil {
		t.Fatal(err)
	}
}
