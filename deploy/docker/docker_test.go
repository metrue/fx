package docker

import (
	"context"
	"testing"
	"time"
)

func TestDocker(t *testing.T) {
	ctx := context.Background()
	cli, err := CreateClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	workdir := "./fixture"
	name := "helloworld"
	ports := []int32{12345, 12346}
	if err := cli.Deploy(ctx, workdir, name, ports); err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	if err := cli.Destroy(ctx, name); err != nil {
		t.Fatal(err)
	}
}
