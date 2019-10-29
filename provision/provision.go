package provision

import (
	"fmt"
	"os"
	"sync"

	"github.com/apex/log"
	"github.com/metrue/fx/config"
	"github.com/metrue/fx/constants"
	"github.com/metrue/fx/pkg/command"
	ssh "github.com/metrue/go-ssh-client"
)

// Provisioner provision
type Provisioner interface {
	Start() error
}

// Provisionor provision-or
type Provisionor struct {
	sshClient ssh.Client

	host config.Host
}

// New new provision
func New(host config.Host) *Provisionor {
	p := &Provisionor{host: host}
	if host.IsRemote() {
		p.sshClient = ssh.New(host.Host).
			WithUser(host.User).
			WithPassword(host.Password)
	}
	return p
}

// NewWithHost create a provisionor with host, user, and password
func NewWithHost(host string, user string, password string) *Provisionor {
	p := &Provisionor{
		host: config.Host{
			Host:     host,
			Password: password,
			User:     user,
		},
	}
	if p.host.IsRemote() {
		p.sshClient = ssh.New(host).
			WithUser(user).
			WithPassword(password)
	}
	return p
}

// IsFxAgentRunning check if fx-agent is running on host
func (p *Provisionor) IsFxAgentRunning() bool {
	script := fmt.Sprintf("docker inspect %s", constants.AgentContainerName)
	var cmd *command.Command
	if p.host.IsRemote() {
		cmd = command.New("inspect fx-agent", script, command.NewRemoteRunner(p.sshClient))
	} else {
		cmd = command.New("inspect fx-agent", script, command.NewLocalRunner())
	}
	output, err := cmd.Exec()
	if os.Getenv("DEBUG") != "" {
		log.Infof(string(output))
	}
	if err != nil {
		return false
	}
	return true
}

// StartFxAgent start fx agent
func (p *Provisionor) StartFxAgent() error {
	script := fmt.Sprintf("docker run -d --name=%s --rm -v /var/run/docker.sock:/var/run/docker.sock -p 0.0.0.0:%s:1234 bobrik/socat TCP-LISTEN:1234,fork UNIX-CONNECT:/var/run/docker.sock", constants.AgentContainerName, constants.AgentPort)
	var cmd *command.Command
	if p.host.IsRemote() {
		cmd = command.New("start fx-agent", script, command.NewRemoteRunner(p.sshClient))
	} else {
		cmd = command.New("start fx-agent", script, command.NewLocalRunner())
	}
	if output, err := cmd.Exec(); err != nil {
		log.Info(string(output))
		return err
	}
	return nil
}

// StopFxAgent stop fx agent
func (p *Provisionor) StopFxAgent() error {
	script := fmt.Sprintf("docker stop %s", constants.AgentContainerName)
	var cmd *command.Command
	if p.host.IsRemote() {
		cmd = command.New("stop fx agent", script, command.NewRemoteRunner(p.sshClient))
	} else {
		cmd = command.New("stop fx agent", script, command.NewLocalRunner())
	}
	if output, err := cmd.Exec(); err != nil {
		log.Infof(string(output))
		return err
	}
	return nil
}

// Start start provision progress
func (p *Provisionor) Start() error {
	scripts := map[string]string{
		"pull java Docker base image":   "docker pull metrue/fx-java-base",
		"pull julia Docker base image":  "docker pull metrue/fx-julia-base",
		"pull python Docker base iamge": "docker pull metrue/fx-python-base",
		"pull node Docker base image":   "docker pull metrue/fx-node-base",
		"pull d Docker base image":      "docker pull metrue/fx-d-base",
		"pull go Docker base image":     "docker pull metrue/fx-go-base",
	}

	var wg sync.WaitGroup
	for n, s := range scripts {
		wg.Add(1)
		go func(name, script string) {
			var cmd *command.Command
			if p.host.IsRemote() {
				cmd = command.New(name, script, command.NewRemoteRunner(p.sshClient))
			} else {
				cmd = command.New(name, script, command.NewLocalRunner())
			}
			if _, err := cmd.Exec(); err != nil {
				log.Fatalf("Provision:%s: %s", cmd.Name, err)
			} else {
				log.Infof("Provision:%s: \u2713", cmd.Name)
			}
			wg.Done()
		}(n, s)
	}
	wg.Wait()
	return nil
}
