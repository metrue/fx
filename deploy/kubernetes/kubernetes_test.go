package kubernetes

import (
	"context"
	"os"
	"testing"

	"github.com/metrue/fx/types"
)

func TestK8SRunner(t *testing.T) {
	workdir := "./fixture"
	name := "hello"
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
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		t.Skip("skip test since no KUBECONFIG given in environment variable")
	}
	k8s, err := Create()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	if err := k8s.Deploy(ctx, workdir, name, bindings); err != nil {
		t.Fatal(err)
	}

	if err := k8s.Destroy(ctx, name); err != nil {
		t.Fatal(err)
	}
}
