package docker

import (
	"context"
	"testing"
	"time"

	"github.com/metrue/fx/types"
)

func TestDocker(t *testing.T) {
	ctx := context.Background()
	cli, err := CreateClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	workdir := "./fixture"
	name := "helloworld"
	bindings := []types.PortBinding{
		types.PortBinding{
			ServiceBindingPort:  80,
			ContainerExposePort: 3000,
		},
		types.PortBinding{
			ServiceBindingPort:  443,
			ContainerExposePort: 3000,
		},
	}
	if err := cli.Deploy(ctx, workdir, name, bindings); err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	if err := cli.Destroy(ctx, name); err != nil {
		t.Fatal(err)
	}
}
