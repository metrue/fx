package k8s

// CreateProvisioner create a provisioner
func CreateProvisioner(master MasterNode, agents []AgentNode) *Provisioner {
	return New(master, agents)
}

// CreateDeployer create a deployer
func CreateDeployer(kubeconfig string) (*K8S, error) {
	return Create(kubeconfig)
}
