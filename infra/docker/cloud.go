package docker

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/metrue/fx/infra"
	"github.com/metrue/fx/types"
	"github.com/metrue/go-ssh-client"
)

// Cloud define a docker host
type Cloud struct {
	IP      string `json:"ip"`
	User    string `json:"user"`
	Port    string `json:"port"`
	Type    string `json:"type"`
	KeyFile string `json:"key_file"`

	sshClient ssh.Clienter
}

const sshConnectionTimeout = 10 * time.Second

// New new a docker cloud
func New(ip string, user string, port string, keyfile string) *Cloud {
	return &Cloud{
		IP:      ip,
		User:    user,
		Port:    port,
		KeyFile: keyfile,
		Type:    types.CloudTypeDocker,
	}
}

// Create a docker node
func Create(ip string, user string, port string, keyfile string) (*Cloud, error) {
	sshClient := ssh.New(ip).WithUser(user).WithKey(keyfile).WithPort(port)
	return &Cloud{
		IP:      ip,
		User:    user,
		Port:    port,
		Type:    types.CloudTypeDocker,
		KeyFile: keyfile,

		sshClient: sshClient,
	}, nil
}

// Load a docker node from meta
func Load(meta []byte) (*Cloud, error) {
	return nil, nil
}

// Provision a host
func (c *Cloud) Provision() error {
	if err := c.runCmd(infra.Scripts["docker_version"].(string)); err != nil {
		if err := c.runCmd(infra.Scripts["install_docker"].(string)); err != nil {
			return err
		}

		if err := c.runCmd(infra.Scripts["start_dockerd"].(string)); err != nil {
			return err
		}
	}

	if err := c.runCmd(infra.Scripts["check_fx_agent"].(string)); err != nil {
		if err := c.runCmd(infra.Scripts["start_fx_agent"].(string)); err != nil {
			return err
		}
	}
	return nil
}

// GetType cloud type
func (c *Cloud) GetType() string {
	return c.Type
}

// IsHealth check if cloud is in health
func (c *Cloud) IsHealth() (bool, error) {
	local := c.IP == "127.0.0.1" || c.IP == "localhost"
	if !local || os.Getenv("CI") != "" {
		ok, err := c.sshClient.Connectable(sshConnectionTimeout)
		if err != nil {
			return false, fmt.Errorf("could not connect to %s@%s:%s via SSH: '%s'", c.User, c.IP, c.Port, err)
		}
		if !ok {
			return false, fmt.Errorf("could not connect to %s@%s:%s via SSH ", c.User, c.IP, c.Port)
		}
	}

	if err := c.runCmd(infra.Scripts["check_fx_agent"].(string)); err != nil {
		if err := c.runCmd(infra.Scripts["start_fx_agent"].(string)); err != nil {
			return false, err
		}
	}
	return true, nil
}

// NOTE only using for unit testing
func (c *Cloud) setsshClient(client ssh.Clienter) {
	c.sshClient = client
}

// nolint:unparam
func (c *Cloud) runCmd(script string, options ...ssh.CommandOptions) error {
	option := ssh.CommandOptions{
		Timeout: sshConnectionTimeout,
	}
	if len(options) >= 1 {
		option = options[0]
	}

	local := c.IP == "127.0.0.1" || c.IP == "localhost"
	if local && os.Getenv("CI") == "" {
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

	return c.sshClient.RunCommand(script, option)
}

var (
	_ infra.Clouder = &Cloud{}
)
