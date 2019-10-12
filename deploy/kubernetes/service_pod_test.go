package kubernetes

import (
	"os"
	"reflect"
	"testing"
)

func TestK8S(t *testing.T) {
	namespace := "default"
	// TODO image is ready on hub.docker.com
	image := "metrue/kube-hello"
	ports := []int32{32300}
	podName := "test-fx-pod"
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		t.Skip("skip test since no KUBECONFIG given in environment variable")
	}
	k8s, err := Create()
	if err != nil {
		t.Fatal(err)
	}

	labels := map[string]string{
		"fx-app": "fx-app",
	}
	newPod, err := k8s.CreatePod(namespace, podName, image, labels)
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

	serviceName := podName + "-svc"
	if _, err := k8s.GetService(namespace, serviceName); err == nil {
		t.Fatalf("should get no service name %s", serviceName)
	}

	svc, err := k8s.CreateService(namespace, serviceName, "NodePort", ports, labels)
	if err != nil {
		t.Fatal(err)
	}
	if svc.Name != serviceName {
		t.Fatalf("should get %s but got %s", serviceName, svc.Name)
	}
	svc, err = k8s.GetService(namespace, serviceName)
	if err != nil {
		t.Fatal(err)
	}
	if svc.Name != serviceName {
		t.Fatalf("should get %s but got %v", serviceName, svc.Name)
	}

	selector := map[string]string{"hello": "world"}
	svc, err = k8s.UpdateService(namespace, serviceName, "NodePort", ports, selector)
	if err != nil {
		t.Fatal(err)
	}
	if svc.Name != serviceName {
		t.Fatalf("should get %s but got %v", serviceName, svc.Name)
	}
	if !reflect.DeepEqual(svc.Spec.Selector, selector) {
		t.Fatalf("should get %v but got %v", selector, svc.Spec.Selector)
	}

	// TODO check service status
	if err := k8s.DeleteService(namespace, serviceName); err != nil {
		t.Fatal(err)
	}
	if err := k8s.DeletePod(namespace, podName); err != nil {
		t.Fatal(err)
	}
}
