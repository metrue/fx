package kubernetes

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// K8S client
type K8S struct {
	*kubernetes.Clientset
}

// GetPod get a pod
func (k *K8S) GetPod(namespace string, name string) (*v1.Pod, error) {
	pod, err := k.CoreV1().Pods(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return pod, nil
}

// ListPods list pods
func (k *K8S) ListPods() (*v1.PodList, error) {
	pods, err := k.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return pods, nil
}

// CreatePod create a pod
func (k *K8S) CreatePod(
	namespace string,
	name string,
	image string,
	port int32,
	labels map[string]string,
) (*v1.Pod, error) {
	container := v1.Container{
		Name:  "fx-placeholder-container-name",
		Image: image,
		Ports: []v1.ContainerPort{
			v1.ContainerPort{
				Name:          "container-fx",
				HostPort:      port,
				ContainerPort: port,
			},
		},
	}

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   name,
			Labels: labels,
		},
		Spec: v1.PodSpec{
			Containers:    []v1.Container{container},
			RestartPolicy: "",
			// NodeName string `json:"nodeName,omitempty" protobuf:"bytes,10,opt,name=nodeName"`
		},
	}

	createdPod, err := k.CoreV1().Pods(namespace).Create(pod)
	if err != nil {
		return nil, err
	}
	return createdPod, nil
}

// DeletePod delete a pod
func (k *K8S) DeletePod(namespace string, name string) error {
	// TODO figure how to delete a pod in a elegant way
	options := metav1.DeleteOptions{}
	return k.CoreV1().Pods(namespace).Delete(name, &options)
}

// New create a k8s cluster client
func New(kubeconfig string) (*K8S, error) {
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
