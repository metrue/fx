package kubernetes

import (
	"os"
	"testing"

	"github.com/metrue/fx/types"
)

func TestDeployment(t *testing.T) {
	namespace := "default"
	name := "fx-hello-world"
	image := "metrue/kube-hello"
	selector := map[string]string{
		"app": "fx-app",
	}

	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		t.Skip("skip test since no KUBECONFIG given in environment variable")
	}

	k8s, err := Create()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := k8s.GetDeployment(namespace, name); err == nil {
		t.Fatalf("should get not found error")
	}

	replicas := int32(2)
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
	deployment, err := k8s.CreateDeployment(namespace, name, image, bindings, replicas, selector)
	if err != nil {
		t.Fatal(err)
	}
	if deployment == nil {
		t.Fatalf("deploymetn should not be %v", nil)
	}

	if deployment.Name != name {
		t.Fatalf("should get %s but got %s", name, deployment.Name)
	}

	if *deployment.Spec.Replicas != replicas {
		t.Fatalf("should get %v but got %v", replicas, deployment.Spec.Replicas)
	}

	if err := k8s.DeleteDeployment(namespace, name); err != nil {
		t.Fatal(err)
	}
}
