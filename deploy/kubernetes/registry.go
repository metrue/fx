package kubernetes

import "context"

const name = "fx-docker-registry"
const image = "registry:2"

// SetupRegistry set a registry
func (k *K8S) SetupRegistry(ctx context.Context, namespace string) error {
	// registry exposes port 5000
	selector := map[string]string{"app": "fx-docker-registry"}
	if _, err := k.CreateDeployment(namespace, name, image, []int32{5000}, 1, selector); err != nil {
		return err
	}

	if _, err := k.CreateService(namespace, name, "NodePort", nil, selector); err != nil {
		return err
	}

	return nil
}
