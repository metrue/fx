package env

import (
	"os/exec"
	"sync"

	"github.com/apex/log"
)

// DockerRemoteAPIEndpoint docker remote api
const DockerRemoteAPIEndpoint = "127.0.0.1:1234"

// Init init a host to make fx runnable
func Init() error {
	// enable docker remote api
	// docker run -d -v /var/run/docker.sock:/var/run/docker.sock -p 127.0.0.1:1234:1234 bobrik/socat TCP-LISTEN:1234,fork UNIX-CONNECT:/var/run/docker.sock
	cmd := exec.Command(
		"docker",
		"run",
		"-d",
		"-v",
		"/var/run/docker.sock:/var/run/docker.sock",
		"-p",
		DockerRemoteAPIEndpoint+":1234",
		"bobrik/socat",
		"TCP-LISTEN:1234,fork",
		"UNIX-CONNECT:/var/run/docker.sock",
	)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Initialize Environment: %v", err)
		return err
	}
	log.Infof("Initialize Environment: \u2713 %s", stdoutStderr)

	baseImages := []string{
		"metrue/fx-java-base",
		"metrue/fx-julia-base",
		"metrue/fx-python-base",
		"metrue/fx-node-base",
		"metrue/fx-d-base",
	}

	var wg sync.WaitGroup
	for _, image := range baseImages {
		wg.Add(1)
		go func(img string) {
			cmd := exec.Command(
				"docker",
				"pull",
				img,
			)
			stdoutStderr, err := cmd.CombinedOutput()
			if err != nil {
				log.Fatalf("Pulling %s failed: %v", img, stdoutStderr)
			} else {
				log.Infof("Pull %s ok", img)
			}
			wg.Done()
		}(image)
	}
	wg.Wait()

	return nil
}
