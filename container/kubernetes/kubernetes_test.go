package kubernetes

import (
	"os"
	"testing"
)

func TestK8SRunner(t *testing.T) {
	// TODO image is ready on hub.docker.com
	name := "fx-test-func"
	image := "metrue/kube-hello"
	port := int32(3000)
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		t.Skip("skip test since no KUBECONFIG given in environment variable")
	}
	k8s, err := Create()
	if err != nil {
		t.Fatal(err)
	}

	var svc []byte
	if err := k8s.Deploy(name, image, port, &svc); err != nil {
		t.Fatal(err)
	}

	if err := k8s.Destroy(name, svc); err != nil {
		t.Fatal(err)
	}
}
