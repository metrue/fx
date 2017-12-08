package env

import api "github.com/metrue/fx/server/docker-api"

func PullBaseDockerImage(verbose bool) {
	baseImages := []string{
		"metrue/fx-java-base",
		"metrue/fx-julia-base",
		"metrue/fx-python-base",
		"metrue/fx-node-base",
	}

	task := func(image string, verbose bool) {
		api.Pull(image, verbose)
	}

	for _, image := range baseImages {
		go task(image, verbose)
	}
}

// Init creates the server
func Init(verbose bool) {
	PullBaseDockerImage(verbose)
}
