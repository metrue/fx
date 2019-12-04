package k3s

import (
	"context"
	"fmt"
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
func Create(kubeconfig string) (*K3S, error) {
	if os.Getenv("KUBECONFIG") != "" {
		kubeconfig = os.Getenv("KUBECONFIG")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
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
func (k *K3S) GetStatus(ctx context.Context, name string) (types.Service, error) {
	svc, err := k.GetService(namespace, name)
	service := types.Service{}
	if err != nil {
		return service, err
	}

	service.Host = svc.Spec.ClusterIP
	if len(svc.Spec.ExternalIPs) > 0 {
		service.Host = svc.Spec.ExternalIPs[0]
	}

	for _, port := range svc.Spec.Ports {
		// TODO should clearify which port (target port, node port) should use
		service.Port = int(port.Port)
		break
	}
	return service, nil
}

// List services
func (k *K3S) List(ctx context.Context, name string) ([]types.Service, error) {
	return []types.Service{}, nil
}

// Ping health check of infra
func (k *K3S) Ping(ctx context.Context) error {
	// Does not find any ping method for k8s
	nodes, err := k.ListNodes()
	if err != nil {
		return err
	}
	if len(nodes.Items) <= 0 {
		return fmt.Errorf("no available nodes")
	}
	return nil
}

var (
	_ deploy.Deployer = &K3S{}
)
