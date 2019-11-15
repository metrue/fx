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

	fn := types.Func{
		Language: "node",
		Source: `
module.exports = (ctx) => {
	ctx.body = 'hello world'
}
`,
	}
	if err := cli.Deploy(ctx, fn, name, bindings); err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	if err := cli.Destroy(ctx, name); err != nil {
		t.Fatal(err)
	}
}
