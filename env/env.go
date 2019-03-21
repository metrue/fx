package env

import (
	"encoding/json"
	"os/exec"
	"sync"

	"github.com/apex/log"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

// DockerRemoteAPIEndpoint docker remote api
const DockerRemoteAPIEndpoint = "127.0.0.1:1234"

type containerInfo struct {
	ID         string                     `json:"Id"`
	State      dockerTypes.ContainerState `json:"State"`
	Image      string                     `json:"Image"`
	HostConfig container.HostConfig       `json:"HostConfig"`
}

const name = "docker-sock-proxy-for-fx"

func runProxy() error {
	// enable docker remote api
	// docker run -d -v /var/run/docker.sock:/var/run/docker.sock -p 127.0.0.1:1234:1234 bobrik/socat TCP-LISTEN:1234,fork UNIX-CONNECT:/var/run/docker.sock
	cmd := exec.Command(
		"docker",
		"run",
		"--name="+name,
		"-d",
		"-v",
		"/var/run/docker.sock:/var/run/docker.sock",
		"-p",
		DockerRemoteAPIEndpoint+":1234",
		"bobrik/socat",
		"TCP-LISTEN:1234,fork",
		"UNIX-CONNECT:/var/run/docker.sock",
	)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

func proxyDockerSock() error {
	cmd := exec.Command(
		"docker",
		"inspect",
		name,
	)
	var infos []containerInfo
	stdoutStderr, err := cmd.CombinedOutput()
	if err := json.Unmarshal(stdoutStderr, &infos); err != nil {
		// no proxy container created
		return runProxy()
	}

	const containerStateCreated = "created"
	const containerStateRestarting = "restarting"
	const containerStateRunning = "running"
	const containerStatePaused = "paused"
	const containerStateExited = "exited"

	state := infos[0].State.Status
	if state == containerStateRestarting ||
		state == containerStateRunning { // FIXME should wait for it ready
		return nil
	}

	if state == containerStateCreated ||
		state == containerStatePaused ||
		state == containerStateExited {
		cmd := exec.Command(
			"docker",
			"start",
			name,
		)
		_, err = cmd.CombinedOutput()
		if err != nil {
			return err
		}
	}
	return nil
}

// Init init a host to make fx runnable
func Init() error {
	if err := proxyDockerSock(); err != nil {
		log.Fatalf("Proxy Docker Remote API Sock Failed: %v", err)
		return err
	}
	log.Info("Proxy Docker Remote API Sock: \u2713")

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
				log.Fatalf("Pulling %s failed: %v", img, string(stdoutStderr))
			} else {
				log.Infof("%s Pulled: \u2713", img)
			}
			wg.Done()
		}(image)
	}
	wg.Wait()

	return nil
}
