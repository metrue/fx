package kubernetes

import (
	"context"
	"os"
	"testing"
)

func TestK8SRunner(t *testing.T) {
	workdir := "./fixture"
	name := "hello"
	ports := []int32{32300}
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		t.Skip("skip test since no KUBECONFIG given in environment variable")
	}
	k8s, err := Create()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	if err := k8s.Deploy(ctx, workdir, name, ports); err != nil {
		t.Fatal(err)
	}

	if err := k8s.Destroy(ctx, name); err != nil {
		t.Fatal(err)
	}
}
