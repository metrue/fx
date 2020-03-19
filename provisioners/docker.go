package provisioners

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/metrue/go-ssh-client"
)

const sshConnectionTimeout = 10 * time.Second

var scripts = map[string]interface{}{
	"docker_version": "docker version",
	"install_docker": "curl -fsSL https://download.docker.com/linux/static/stable/x86_64/docker-18.06.3-ce.tgz -o docker.tgz && tar zxvf docker.tgz && mv docker/* /usr/bin && rm -rf docker docker.tgz",
	"start_dockerd":  "dockerd >/dev/null 2>&1 & sleep 2",
	"check_fx_agent": "docker inspect fx-agent",
	"start_fx_agent": "docker run -d --name=fx-agent --rm -v /var/run/docker.sock:/var/run/docker.sock -p 0.0.0.0:8866:1234 bobrik/socat TCP-LISTEN:1234,fork UNIX-CONNECT:/var/run/docker.sock",
}

// Docker define a fx docker host
type Docker struct {
	sshClient ssh.Clienter
}

// New a docker provioner
func New(sshClient ssh.Clienter) *Docker {
	return &Docker{
		sshClient: sshClient,
	}
}

// Provision a host
func (d *Docker) Provision(ctx context.Context, isRemote bool) error {
	if err := d.runCmd(scripts["docker_version"].(string), isRemote); err != nil {
		if err := d.runCmd(scripts["install_docker"].(string), isRemote); err != nil {
			return err
		}

		if err := d.runCmd(scripts["start_dockerd"].(string), isRemote); err != nil {
			return err
		}
	}

	if err := d.runCmd(scripts["check_fx_agent"].(string), isRemote); err != nil {
		if err := d.runCmd(scripts["start_fx_agent"].(string), isRemote); err != nil {
			return err
		}
	}
	return nil
}

func (d *Docker) runCmd(script string, isRemote bool, options ...ssh.CommandOptions) error {
	option := ssh.CommandOptions{
		Timeout: sshConnectionTimeout,
	}
	if len(options) >= 1 {
		option = options[0]
	}
	if !isRemote {
		params := strings.Split(script, " ")
		if len(params) == 0 {
			return fmt.Errorf("invalid script: %s", script)
		}
		// nolint
		cmd := exec.Command(params[0], params[1:]...)
		cmd.Stdout = option.Stdout
		cmd.Stderr = option.Stderr
		err := cmd.Run()
		if err != nil {
			return err
		}
		return nil
	}
	ok, err := d.sshClient.Connectable(sshConnectionTimeout)
	if err != nil {
		return fmt.Errorf("could not connect via SSH: '%s'", err)
	}
	if !ok {
		return fmt.Errorf("could not connect via SSH")
	}

	return d.sshClient.RunCommand(script, option)
}

var (
	_ Provisioner = &Docker{}
)
