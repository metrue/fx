package kubernetes

import (
	"fmt"
	"os"
	"testing"
)

func TestK8S(t *testing.T) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		t.Skip("skip test since no KUBECONFIG given in environment variable")
	}
	k8s, err := New(kubeconfig)
	if err != nil {
		t.Fatal(err)
	}

	newPod, err := k8s.CreatePod("default", "test-fx-pod", "metrue/kube-hello", 3000)
	if err != nil {
		t.Fatal(err)
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

	fmt.Println(newPod.Name)
}
