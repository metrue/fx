package kubernetes

import (
	"github.com/metrue/fx/runners"
	"k8s.io/client-go/kubernetes"
)

// K8S client
type K8S struct {
	*kubernetes.Clientset
}

// Deploy a image to be a service
func (k *K8S) Deploy(name string, image string, port int32, svc interface{}) error {
	namespace := "default"
	labels := map[string]string{}
	if _, err := k.CreatePod(
		namespace,
		name,
		image,
		port,
		labels,
	); err != nil {
		return err
	}

	const isOnPublicCloud = false
	typ := "LoadBalencer"
	if !isOnPublicCloud {
		typ = "NodePort"
	}
	podsLabels := map[string]string{}
	if _, err := k.CreateService(
		namespace,
		name,
		typ,
		[]int32{port},
		podsLabels,
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
	_ runners.Runner = &K8S{}
)
