package kubernetes

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	runtime "github.com/metrue/fx/container_runtimes/docker/sdk"
	"github.com/metrue/fx/deploy"
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
func (k *K8S) Deploy(
	ctx context.Context,
	workdir string,
	name string,
	ports []int32,
) error {
	namespace := "default"

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
	labels := map[string]string{
		"fx-app": "fx-app-" + uuid.New().String(),
	}
	if _, err := k.CreatePod(
		namespace,
		name,
		image,
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
		ports,
		labels,
	); err != nil {
		return err
	}
	return nil
}

// Update a service
func (k *K8S) Update(ctx context.Context, name string) error {
	return nil
}

// Destroy a service
func (k *K8S) Destroy(ctx context.Context, name string) error {
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
func (k *K8S) GetStatus(ctx context.Context, name string) error {
	return nil
}

var (
	_ deploy.Deployer = &K8S{}
)
