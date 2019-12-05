package docker

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/metrue/fx/constants"
	"github.com/metrue/fx/infra"
	"github.com/metrue/fx/pkg/spinner"
	sshOperator "github.com/metrue/go-ssh-client"
)

// Provisioner docker host
type Provisioner struct {
	IP   string
	User string
}

// NewProvisioner new a docker object
func NewProvisioner(ip string, user string) *Provisioner {
	return &Provisioner{
		IP:   ip,
		User: user,
	}
}

// Provision provision a host, install docker and start dockerd
func (d *Provisioner) Provision() (config []byte, err error) {
	spinner.Start("provisioning")
	defer func() {
		spinner.Stop("provisioning", err)
	}()

	// TODO clean up, skip check localhost or not if in CICD env
	if os.Getenv("CICD") != "" {
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

	if d.isLocalHost() {
		if !d.hasDocker() {
			return nil, fmt.Errorf("please make sure docker installed and running")
		}

		if err := d.StartFxAgentLocally(); err != nil {
			return nil, err
		}

		config, _ := json.Marshal(map[string]string{
			"ip":   d.IP,
			"user": d.User,
		})
		return config, nil
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
	return json.Marshal(map[string]string{
		"ip":   d.IP,
		"user": d.User,
	})
}

func (d *Provisioner) isLocalHost() bool {
	return strings.ToLower(d.IP) == "localhost" || d.IP == "127.0.0.1"
}

func (d *Provisioner) hasDocker() bool {
	cmd := exec.Command("docker", "version")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// HealthCheck check healthy status of host
func (d *Provisioner) HealthCheck() (bool, error) {
	if d.isLocalHost() {
		return d.IfFxAgentRunningLocally(), nil
	}
	return d.IfFxAgentRunning(), nil
}

// Install docker on host
func (d *Provisioner) Install() error {
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
func (d *Provisioner) StartDockerd() error {
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
func (d *Provisioner) StartFxAgent() error {
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
func (d *Provisioner) StartFxAgentLocally() error {
	startCmd := fmt.Sprintf("docker run -d --name=%s --rm -v /var/run/docker.sock:/var/run/docker.sock -p 0.0.0.0:%s:1234 bobrik/socat TCP-LISTEN:1234,fork UNIX-CONNECT:/var/run/docker.sock", constants.AgentContainerName, constants.AgentPort)
	params := strings.Split(startCmd, " ")
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

// IfFxAgentRunningLocally check if fx agent is running
func (d *Provisioner) IfFxAgentRunningLocally() bool {
	cmd := exec.Command("docker", "inspect", "fx-agent")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

// IfFxAgentRunning check if fx agent is running
func (d *Provisioner) IfFxAgentRunning() bool {
	inspectCmd := infra.Sudo("docker inspect fx-agent", d.User)
	sshKeyFile, _ := infra.GetSSHKeyFile()
	sshPort := infra.GetSSHPort()
	ssh := sshOperator.New(d.IP).WithUser(d.User).WithKey(sshKeyFile).WithPort(sshPort)
	if err := ssh.RunCommand(inspectCmd, sshOperator.CommandOptions{
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
		Stderr: os.Stderr,
	}); err != nil {
		return false
	}
	return true
}

var _ infra.Provisioner = &Provisioner{}
