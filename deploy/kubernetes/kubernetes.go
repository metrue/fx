package kubernetes

import (
	"context"

	runtime "github.com/metrue/fx/container_runtimes/docker/sdk"
	"github.com/metrue/fx/deploy"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// K8S client
type K8S struct {
	*kubernetes.Clientset
}

const namespace = "default"

// Create a k8s cluster client
func Create() (*K8S, error) {
	config, err := clientcmd.BuildConfigFromKubeconfigGetter("", clientcmd.NewDefaultClientConfigLoadingRules().Load)
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
func (k *K8S) Deploy(
	ctx context.Context,
	workdir string,
	name string,
	ports []int32,
) error {
	dockerClient, err := runtime.CreateClient(ctx)
	if err != nil {
		return err
	}
	if err := dockerClient.BuildImage(ctx, workdir, name); err != nil {
		return err
	}
	image, err := dockerClient.PushImage(ctx, name)
	if err != nil {
		return err
	}

	// By using a label selector between Pod and Service, we can link Service and Pod directly, it means a Endpoint will
	// be created automatically, then incoming traffic to Service will be forward to Pod.
	// Then we have no need to create Endpoint manually anymore.
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
			replicas,
			selector,
		); err != nil {
			return err
		}
	} else {
		if _, err := k.UpdateDeployment(namespace, name, image, replicas, selector); err != nil {
			return err
		}
	}

	// TODO fx should be able to know what's the target Kubernetes service platform
	// it's going to deploy to
	const isOnPublicCloud = true
	typ := "LoadBalancer"
	if !isOnPublicCloud {
		typ = "NodePort"
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
func (k *K8S) Update(ctx context.Context, name string) error {
	return nil
}

// Destroy a service
func (k *K8S) Destroy(ctx context.Context, name string) error {
	if err := k.DeleteService(namespace, name); err != nil {
		return err
	}
	if err := k.DeleteDeployment(namespace, name); err != nil {
		return err
	}
	return nil
}

// GetStatus get status of a service
func (k *K8S) GetStatus(ctx context.Context, name string) error {
	return nil
}

var (
	_ deploy.Deployer = &K8S{}
)
