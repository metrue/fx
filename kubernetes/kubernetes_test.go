package kubernetes

import (
	"fmt"
	"testing"
)

func TestK8S(t *testing.T) {
	k8s := New()
	if k8s == nil {
		t.Fatalf("k8s client new failed")
	}

	podList, err := k8s.ListPods()
	if err != nil {
		t.Fatal(err)
	}
	if len(podList.Items) <= 0 {
		t.Fatal("pod number should > 0")
	}

	name := "fx-kube-1-6954f99b9b-ls5zs"
	pod, err := k8s.GetPod("default", name)
	if err != nil {
		t.Fatal(err)
	}
	if pod.Name != name {
		t.Fatalf("should get %s but got %s", name, pod.Name)
	}

	newPod, err := k8s.CreatePod("default", "test-fx-pod", "metrue/kube-hello", 3000)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(newPod.Name)
}
