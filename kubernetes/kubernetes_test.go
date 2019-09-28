package kubernetes

import (
	"os"
	"testing"
)

func TestK8S(t *testing.T) {
	namespace := "default"
	// TODO image is ready on hub.docker.com
	image := "metrue/kube-hello"
	port := int32(3000)
	podName := "test-fx-pod"
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		t.Skip("skip test since no KUBECONFIG given in environment variable")
	}
	k8s, err := New(kubeconfig)
	if err != nil {
		t.Fatal(err)
	}

	newPod, err := k8s.CreatePod(namespace, podName, image, port, map[string]string{})
	if err != nil {
		t.Fatal(err)
	}
	if newPod.Name != podName {
		t.Fatalf("should get %s but got %s", podName, newPod.Name)
	}

	podList, err := k8s.ListPods()
	if err != nil {
		t.Fatal(err)
	}
	if len(podList.Items) <= 0 {
		t.Fatal("pod number should > 0")
	}

	pod := podList.Items[0]
	p, err := k8s.GetPod("default", pod.Name)
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != pod.Name {
		t.Fatalf("should get %s but got %s", pod.Name, p.Name)
	}

	if err := k8s.DeletePod(namespace, podName); err != nil {
		t.Fatal(err)
	}
}
