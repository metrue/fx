package env

import "github.com/metrue/fx/docker-api"

//PullBaseDockerImage fetch base images from the registry
func PullBaseDockerImage(verbose bool) {
	baseImages := []string{
		"metrue/fx-java-base",
		"metrue/fx-julia-base",
		"metrue/fx-python-base",
		"metrue/fx-node-base",
		"metrue/fx-d-base",
	}

	task := func(image string, verbose bool) {
		docker.Pull(image, verbose)
	}

	for _, image := range baseImages {
		go task(image, verbose)
	}
}

// Init creates the server
func Init(verbose bool) {
	PullBaseDockerImage(verbose)
}
