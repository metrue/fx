package k8s

// CreateDeployer create a deployer
func CreateDeployer(kubeconfig string) (*K8S, error) {
	return Create(kubeconfig)
}
