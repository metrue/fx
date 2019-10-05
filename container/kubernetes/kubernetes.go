package kubernetes

import (
	"fmt"
	"os"

	"github.com/metrue/fx/container"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// K8S client
type K8S struct {
	*kubernetes.Clientset
}

// Create a k8s cluster client
func Create() (*K8S, error) {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		return nil, fmt.Errorf("KUBECONFIG not given")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &K8S{clientset}, nil
}

// Deploy a image to be a service
func (k *K8S) Deploy(name string, image string, port int32, svc interface{}) error {
	namespace := "default"
	labels := map[string]string{
		"fx": "fx",
	}
	if _, err := k.CreatePod(
		namespace,
		name,
		image,
		port,
		labels,
	); err != nil {
		return err
	}

	// TODO fx should be able to know what's the target Kubernetes service platform
	// it's going to deploy to
	const isOnPublicCloud = true
	typ := "LoadBalancer"
	if !isOnPublicCloud {
		typ = "NodePort"
	}
	if _, err := k.CreateService(
		namespace,
		name,
		typ,
		[]int32{port},
		labels,
	); err != nil {
		return err
	}
	return nil
}

// Update a service
func (k *K8S) Update(name string, svc interface{}) error {
	return nil
}

// Destroy a service
func (k *K8S) Destroy(name string, svc interface{}) error {
	const namespace = "default"
	if err := k.DeleteService(namespace, name); err != nil {
		return err
	}
	if err := k.DeletePod(namespace, name); err != nil {
		return err
	}
	return nil
}

// GetStatus get status of a service
func (k *K8S) GetStatus(name string, svc interface{}) error {
	return nil
}

var (
	_ container.Runner = &K8S{}
)
