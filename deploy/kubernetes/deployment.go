package kubernetes

import (
	"fmt"

	"github.com/metrue/fx/types"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func generateDeploymentSpec(
	name string,
	image string,
	bindPorts []types.PortBinding,
	replicas int32,
	selector map[string]string,
) *appsv1.Deployment {
	ports := []apiv1.ContainerPort{}
	for index, binding := range bindPorts {
		ports = append(ports, apiv1.ContainerPort{
			Name:          fmt.Sprintf("fx-container-%d", index),
			ContainerPort: binding.ContainerExposePort,
		})
	}

	container := apiv1.Container{
		Name:            "fx-placeholder-container-name",
		Image:           image,
		Ports:           ports,
		ImagePullPolicy: v1.PullNever,
	}
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: selector,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: selector,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{container},
				},
			},
		},
	}
}

// GetDeployment get a deployment
func (k *K8S) GetDeployment(namespace string, name string) (*appsv1.Deployment, error) {
	return k.AppsV1().Deployments(namespace).Get(name, metav1.GetOptions{})
}

// CreateDeployment create a deployment
func (k *K8S) CreateDeployment(
	namespace string,
	name string,
	image string,
	ports []types.PortBinding,
	replicas int32,
	selector map[string]string,
) (*appsv1.Deployment, error) {
	deployment := generateDeploymentSpec(name, image, ports, replicas, selector)
	return k.AppsV1().Deployments(namespace).Create(deployment)
}

// UpdateDeployment update a deployment
func (k *K8S) UpdateDeployment(
	namespace string,
	name string,
	image string,
	ports []types.PortBinding,
	replicas int32,
	selector map[string]string,
) (*appsv1.Deployment, error) {
	deployment := generateDeploymentSpec(name, image, ports, replicas, selector)
	return k.AppsV1().Deployments(namespace).Update(deployment)
}

// DeleteDeployment delete a deployment
func (k *K8S) DeleteDeployment(namespace string, name string) error {
	return k.AppsV1().Deployments(namespace).Delete(name, &metav1.DeleteOptions{})
}

// CreateDeploymentWithInitContainer create a deployment which will wait InitContainer to do the image build before function container start
func (k *K8S) CreateDeploymentWithInitContainer(
	namespace string,
	name string,
	ports []types.PortBinding,
	replicas int32,
	selector map[string]string,
) (*appsv1.Deployment, error) {
	deployment := generateDeploymentSpec(name, name, ports, replicas, selector)
	updatedDeployment := injectInitContainer(name, deployment)
	fmt.Println(updatedDeployment)
	return k.AppsV1().Deployments(namespace).Create(updatedDeployment)
}
