package kubernetes

import (
	"context"
	"os"
	"testing"
)

func TestK8SDeployer(t *testing.T) {
	workdir := "./fixture"
	name := "hello"
	ports := []int32{32300}
	kubeconfig := os.Getenv("KUBECONFIG")
	username := os.Getenv("DOCKER_USERNAME")
	password := os.Getenv("DOCKER_PASSWORD")
	if kubeconfig == "" || username == "" || password == "" {
		t.Skip("skip test since no KUBECONFIG, DOCKER_USERNAME and DOCKER_PASSWORD given in environment variable")
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
