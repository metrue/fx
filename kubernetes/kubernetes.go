package kubernetes

import (
	"k8s.io/client-go/kubernetes"
)

// K8S client
type K8S struct {
	*kubernetes.Clientset
}
