package api

import (
	"os/exec"

	"github.com/apex/log"
)

// Init init a host to be a fx running
func (api *API) Init() error {
	// enable docker remote api
	// docker run -d -v /var/run/docker.sock:/var/run/docker.sock -p 127.0.0.1:1234:1234 bobrik/socat TCP-LISTEN:1234,fork UNIX-CONNECT:/var/run/docker.sock
	cmd := exec.Command(
		"docker",
		"run",
		"-d",
		"-v",
		"/var/run/docker.sock:/var/run/docker.sock",
		"-p",
		"127.0.0.1:1234:1234",
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
	return nil
}
