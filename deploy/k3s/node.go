package k3s

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ListNodes list node
func (k *K3S) ListNodes() (*v1.NodeList, error) {
	nodes, err := k.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return nodes, nil
}
