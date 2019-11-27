package config

// DockerCloud defines a host with bare docker
type DockerCloud struct {
	Host string `json:"host"`
	User string `json:"user"`
}

// K8SCloud defines a k8s cluster with kubeconfig
type K8SCloud struct {
	KubeConfig string `json:"kubeconfig"`
}
