package kubernetes

import (
	"os"
	"testing"
)

func TestConfigMap(t *testing.T) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		t.Skip("skip test since no KUBECONFIG given in environment variable")
	}

	k8s, err := Create()
	if err != nil {
		t.Fatal(err)
	}

	namespace := "default"
	name := "test-configmap"
	data := map[string]string{
		"message": "hello world",
	}
	cm, err := k8s.CreateConfigMap(namespace, name, data)
	if err != nil {
		t.Fatal(err)
	}
	if cm.Name != name {
		t.Fatalf("should get %s but got %s", name, cm.Name)
	}

	if err != k8s.DeleteConfigMap(namespace, name) {
		t.Fatal(err)
	}
}
