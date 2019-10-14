package kubernetes

import (
	"context"

	"github.com/metrue/fx/types"
)

const name = "fx-docker-registry"
const image = "registry:2"

// SetupRegistry set a registry
func (k *K8S) SetupRegistry(ctx context.Context, namespace string) error {
	// registry exposes port 5000
	selector := map[string]string{"app": "fx-docker-registry"}
	bindings := []types.PortBinding{
		types.PortBinding{
			ServiceBindingPort:  80,
			ContainerExposePort: 5000,
		},
		types.PortBinding{
			ServiceBindingPort:  443,
			ContainerExposePort: 5000,
		},
	}
	if _, err := k.CreateDeployment(namespace, name, image, bindings, 1, selector); err != nil {
		return err
	}

	if _, err := k.CreateService(namespace, name, "NodePort", bindings, selector); err != nil {
		return err
	}

	return nil
}
