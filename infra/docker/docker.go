package docker

import (
	"fmt"

	"github.com/metrue/fx/constants"
	sshOperator "github.com/metrue/go-ssh-client"
	"github.com/mitchellh/go-homedir"
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

// Install docker on host
func (d *Docker) Install() error {
	installCmd := "curl -fsSL https://download.docker.com/linux/static/stable/x86_64/docker-18.06.3-ce.tgz -o docker.tgz && tar zxvf docker.tgz && sudo mv docker/* /usr/bin && rm -rf docker docker.tgz"
	publicKey, err := homedir.Expand("~/.ssh/id_rsa")
	if err != nil {
		return err
	}
	ssh := sshOperator.New(d.IP).WithUser(d.User).WithKey(publicKey)
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
	installCmd := "sudo dockerd >/dev/null 2>&1 & sleep 2"
	publicKey, err := homedir.Expand("~/.ssh/id_rsa")
	if err != nil {
		return err
	}
	ssh := sshOperator.New(d.IP).WithUser(d.User).WithKey(publicKey)
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
	startCmd := fmt.Sprintf("docker run -d --name=%s --rm -v /var/run/docker.sock:/var/run/docker.sock -p 0.0.0.0:%s:1234 bobrik/socat TCP-LISTEN:1234,fork UNIX-CONNECT:/var/run/docker.sock", constants.AgentContainerName, constants.AgentPort)
	publicKey, err := homedir.Expand("~/.ssh/id_rsa")
	if err != nil {
		return err
	}
	ssh := sshOperator.New(d.IP).WithUser(d.User).WithKey(publicKey)
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
