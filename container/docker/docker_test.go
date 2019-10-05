package docker

import (
	"context"
	"testing"
	"time"
)

func TestDocker(t *testing.T) {
	cli, err := CreateClient()
	if err != nil {
		t.Fatal(err)
	}

	name := "hello-fx"
	image := "nginxdemos/hello"
	ports := []int32{12345, 12346}
	if err := cli.Deploy(context.Background(), name, image, ports); err != nil {
		t.Fatal(err)
	}

	time.Sleep(1 * time.Second)

	if err := cli.Destroy(context.Background(), name); err != nil {
		t.Fatal(err)
	}
}
