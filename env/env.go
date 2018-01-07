package env

import (
	docker "github.com/metrue/fx/docker-api"
	"github.com/pkg/errors"
)

type PullTask struct {
	ImageName string
	Err       error
}

func newPullTask(imageName string, result error) PullTask {
	return PullTask{
		ImageName: imageName,
		Err:       result,
	}
}

//PullBaseDockerImage fetch base images from the registry
func PullBaseDockerImage(verbose bool) []PullTask {
	baseImages := []string{
		"metrue/fx-java-base",
		"metrue/fx-julia-base",
		"metrue/fx-python-base",
		"metrue/fx-node-base",
		"metrue/fx-d-base",
	}

	count := len(baseImages)
	results := make(chan PullTask, count)

	task := func(image string, verbose bool) error {
		return docker.Pull(image, verbose)
	}

	for _, image := range baseImages {
		go func(img string) {
			err := task(img, verbose)
			results <- newPullTask(img, err)
		}(image)
	}

	var pullResutls []PullTask
	for result := range results {
		pullResutls = append(pullResutls, result)

		if len(pullResutls) == count {
			close(results)
		}
	}

	return pullResutls
}

// Init creates the server
func Init(verbose bool) (ret []PullTask, err error) {
	err = docker.Info()
	if err != nil {
		err = errors.Wrap(err, "docker info")
	} else {
		ret = PullBaseDockerImage(verbose)
	}

	return ret, err
}
