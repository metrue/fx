package kubernetes

import (
	"context"
	"os"
	"testing"
)

func TestSetupRegistry(t *testing.T) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		t.Skip("skip test since no KUBECONFIG given in environment variable")
	}
	k8s, err := Create()
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	namespace := "default"
	if err := k8s.SetupRegistry(ctx, namespace); err != nil {
		t.Fatal(err)
	}
}
