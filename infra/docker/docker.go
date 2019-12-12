package docker

import containerruntimes "github.com/metrue/fx/container_runtimes"

// CreateDeployer create a deployer
func CreateDeployer(client containerruntimes.ContainerRuntime) (*Deployer, error) {
	return &Deployer{cli: client}, nil
}
