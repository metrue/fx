package docker

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
	fmt.Println("+++++")
	// TODO clean up, skip check localhost or not if in CICD env
	if d.isLocalHost() {
		fmt.Println("+++++")
		if os.Getenv("CICD") == "" {
			fmt.Println("+++++")
			if !d.hasDocker() {
				return nil, fmt.Errorf("please make sure docker installed and running")
			}
			fmt.Println("+++++")

			if err := d.StartFxAgentLocally(); err != nil {
				return nil, err
			}

			fmt.Println("+++++")
			config, _ := json.Marshal(map[string]string{
				"ip":   d.IP,
				"user": d.User,
			})
			return config, nil
		}
	}

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

func (d *Docker) isLocalHost() bool {
	return strings.ToLower(d.IP) == "localhost" || d.IP == "127.0.0.1"
}

func (d *Docker) hasDocker() bool {
	cmd := exec.Command("docker", "version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
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
	if err := ssh.RunCommand(installCmd, sshOperator.CommandOptions{
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
		Stderr: os.Stderr,
	}); err != nil {
		fmt.Println("install docker failed \n================")
		fmt.Println(err)
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
	installCmd := fmt.Sprintf("%s dockerd >/dev/null 2>&1 & sleep 2", sudo)
	sshKeyFile, _ := infra.GetSSHKeyFile()
	sshPort := infra.GetSSHPort()
	ssh := sshOperator.New(d.IP).WithUser(d.User).WithKey(sshKeyFile).WithPort(sshPort)
	if err := ssh.RunCommand(installCmd, sshOperator.CommandOptions{
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
		Stderr: os.Stderr,
	}); err != nil {
		fmt.Println("start dockerd failed \n================")
		fmt.Println(err)
		fmt.Println("================")
		return err
	}

	return nil
}

// StartFxAgent start fx agent
func (d *Docker) StartFxAgent() error {
	startCmd := fmt.Sprintf("sleep 3 && docker stop %s || true && docker run -d --name=%s --rm -v /var/run/docker.sock:/var/run/docker.sock -p 0.0.0.0:%s:1234 bobrik/socat TCP-LISTEN:1234,fork UNIX-CONNECT:/var/run/docker.sock", constants.AgentContainerName, constants.AgentContainerName, constants.AgentPort)
	sshKeyFile, _ := infra.GetSSHKeyFile()
	sshPort := infra.GetSSHPort()
	ssh := sshOperator.New(d.IP).WithUser(d.User).WithKey(sshKeyFile).WithPort(sshPort)
	if err := ssh.RunCommand(startCmd, sshOperator.CommandOptions{
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
		Stderr: os.Stderr,
	}); err != nil {
		fmt.Println("start fx agent failed \n================")
		fmt.Println(err)
		fmt.Println("================")
		return err
	}
	return nil
}

// StartFxAgentLocally start fx agent
func (d *Docker) StartFxAgentLocally() error {
	startCmd := fmt.Sprintf("docker run -d --name=%s --rm -v /var/run/docker.sock:/var/run/docker.sock -p 0.0.0.0:%s:1234 bobrik/socat TCP-LISTEN:1234,fork UNIX-CONNECT:/var/run/docker.sock", constants.AgentContainerName, constants.AgentPort)
	params := strings.Split(startCmd, " ")
	fmt.Println(params)
	var cmd *exec.Cmd
	if len(params) > 1 {
		// nolint: gosec
		cmd = exec.Command(params[0], params[1:]...)
	} else {
		// nolint: gosec
		cmd = exec.Command(params[0])
	}
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Println(string(out))
		return err
	}
	return nil
}

var _ infra.Infra = &Docker{}
