package docker

import (
	"encoding/json"
	"fmt"

	"github.com/metrue/fx/constants"
	"github.com/metrue/fx/infra"
	sshOperator "github.com/metrue/go-ssh-client"
)

// Docker docker host
type Docker struct {
	IP   string
	User string
}

// New new a docker object
func New(ip string, user string) *Docker {
	return &Docker{
		IP:   ip,
		User: user,
	}
}

// Provision provision a host, install docker and start dockerd
func (d *Docker) Provision() ([]byte, error) {
	if err := d.Install(); err != nil {
		return nil, err
	}
	if err := d.StartDockerd(); err != nil {
		return nil, err
	}
	if err := d.StartFxAgent(); err != nil {
		return nil, err
	}
	config, _ := json.Marshal(map[string]string{
		"ip":   d.IP,
		"user": d.User,
	})
	return config, nil
}

// HealthCheck check healthy status of host
func (d *Docker) HealthCheck() (bool, error) {
	// TODO
	return true, nil
}

// Install docker on host
func (d *Docker) Install() error {
	sudo := ""
	if d.User != "root" {
		sudo = "sudo"
	}
	installCmd := fmt.Sprintf("curl -fsSL https://download.docker.com/linux/static/stable/x86_64/docker-18.06.3-ce.tgz -o docker.tgz && tar zxvf docker.tgz && %s mv docker/* /usr/bin && rm -rf docker docker.tgz", sudo)
	sshKeyFile, _ := infra.GetSSHKeyFile()
	sshPort := infra.GetSSHPort()
	ssh := sshOperator.New(d.IP).WithUser(d.User).WithKey(sshKeyFile).WithPort(sshPort)
	stdout, stderr, err := ssh.RunCommand(installCmd)
	if err != nil {
		fmt.Println("install docker failed \n================")
		fmt.Println("failed: ", string(stderr))
		fmt.Println("output: ", string(stdout))
		fmt.Println("================")
		return err
	}

	return nil
}

// StartDockerd start dockerd
func (d *Docker) StartDockerd() error {
	sudo := ""
	if d.User != "root" {
		sudo = "sudo"
	}
	installCmd := fmt.Sprintf("%s, dockerd >/dev/null 2>&1 & sleep 2", sudo)
	sshKeyFile, _ := infra.GetSSHKeyFile()
	sshPort := infra.GetSSHPort()
	ssh := sshOperator.New(d.IP).WithUser(d.User).WithKey(sshKeyFile).WithPort(sshPort)
	stdout, stderr, err := ssh.RunCommand(installCmd)
	if err != nil {
		fmt.Println("start dockerd failed \n================")
		fmt.Println("failed: ", string(stderr))
		fmt.Println("output: ", string(stdout))
		fmt.Println("================")
		return err
	}

	return nil
}

// StartFxAgent start fx agent
func (d *Docker) StartFxAgent() error {
	startCmd := fmt.Sprintf("sleep 3 && docker run -d --name=%s --rm -v /var/run/docker.sock:/var/run/docker.sock -p 0.0.0.0:%s:1234 bobrik/socat TCP-LISTEN:1234,fork UNIX-CONNECT:/var/run/docker.sock", constants.AgentContainerName, constants.AgentPort)
	sshKeyFile, _ := infra.GetSSHKeyFile()
	sshPort := infra.GetSSHPort()
	ssh := sshOperator.New(d.IP).WithUser(d.User).WithKey(sshKeyFile).WithPort(sshPort)
	stdout, stderr, err := ssh.RunCommand(startCmd)
	if err != nil {
		fmt.Println("start fx agent failed \n================")
		fmt.Println("failed: ", string(stderr))
		fmt.Println("output: ", string(stdout))
		fmt.Println("================")
		return err
	}
	return nil
}

var _ infra.Infra = &Docker{}
