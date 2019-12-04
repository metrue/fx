package docker

import containerruntimes "github.com/metrue/fx/container_runtimes"

// CreateProvisioner create a provisioner
func CreateProvisioner(ip string, user string) *Provisioner {
	return NewProvisioner(ip, user)
}

// CreateDeployer create a deployer
func CreateDeployer(client containerruntimes.ContainerRuntime) (*Deployer, error) {
	return &Deployer{cli: client}, nil
}
