package k3s

import (
	"context"
	"os"

	"github.com/metrue/fx/deploy"
	"github.com/metrue/fx/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// K3S client
type K3S struct {
	*kubernetes.Clientset
}

const namespace = "default"

// Create a k8s cluster client
func Create() (*K3S, error) {
	config, err := clientcmd.BuildConfigFromKubeconfigGetter("", clientcmd.NewDefaultClientConfigLoadingRules().Load)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return &K3S{clientset}, nil
}

// Deploy a image to be a service
func (k *K3S) Deploy(
	ctx context.Context,
	fn types.Func,
	name string,
	image string,
	ports []types.PortBinding,
) error {
	selector := map[string]string{
		"app": "fx-app-" + name,
	}

	const replicas = int32(3)
	if _, err := k.GetDeployment(namespace, name); err != nil {
		// TODO enable passing replica from fx CLI
		if _, err := k.CreateDeployment(
			namespace,
			name,
			image,
			ports,
			replicas,
			selector,
		); err != nil {
			return err
		}
	} else {
		if _, err := k.UpdateDeployment(
			namespace,
			name,
			image,
			ports,
			replicas,
			selector,
		); err != nil {
			return err
		}
	}

	// TODO fx should be able to know what's the target Kubernetes service platform
	// it's going to deploy to
	typ := "LoadBalancer"
	if os.Getenv("SERVICE_TYPE") != "" {
		typ = os.Getenv("SERVICE_TYPE")
	}

	if _, err := k.GetService(namespace, name); err != nil {
		if _, err := k.CreateService(
			namespace,
			name,
			typ,
			ports,
			selector,
		); err != nil {
			return err
		}
	} else {
		if _, err := k.UpdateService(
			namespace,
			name,
			typ,
			ports,
			selector,
		); err != nil {
			return err
		}
	}
	return nil
}

// Update a service
func (k *K3S) Update(ctx context.Context, name string) error {
	return nil
}

// Destroy a service
func (k *K3S) Destroy(ctx context.Context, name string) error {
	if err := k.DeleteService(namespace, name); err != nil {
		return err
	}
	if err := k.DeleteDeployment(namespace, name); err != nil {
		return err
	}
	return nil
}

// GetStatus get status of a service
func (k *K3S) GetStatus(ctx context.Context, name string) error {
	return nil
}

// List services
func (k *K3S) List(ctx context.Context, name string) ([]types.Service, error) {
	return []types.Service{}, nil
}

var (
	_ deploy.Deployer = &K3S{}
)
