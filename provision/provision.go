package provision

import (
	"fmt"
	"strings"
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

// Start start provision progress
func (p *Provisionor) Start() error {
	startFxAgent := fmt.Sprintf("docker run -d --name=%s --rm -v /var/run/docker.sock:/var/run/docker.sock -p 0.0.0.0:%s:1234 bobrik/socat TCP-LISTEN:1234,fork UNIX-CONNECT:/var/run/docker.sock", constants.AgentContainerName, constants.AgentPort)
	stopFxAgent := fmt.Sprintf("docker stop %s", constants.AgentContainerName)
	scripts := map[string]string{
		"pull java Docker base image":   "docker pull metrue/fx-java-base",
		"pull julia Docker base image":  "docker pull metrue/fx-julia-base",
		"pull python Docker base iamge": "docker pull metrue/fx-python-base",
		"pull node Docker base image":   "docker pull metrue/fx-node-base",
		"pull d Docker base image":      "docker pull metrue/fx-d-base",
		"pull go Docker base image":     "docker pull metrue/fx-go-base",
	}

	agentStartupCmds := []*command.Command{}
	if p.host.IsRemote() {
		agentStartupCmds = append(agentStartupCmds,
			command.New("stop current fx agent", stopFxAgent, command.NewRemoteRunner(p.sshClient)),
			command.New("start fx agent", startFxAgent, command.NewRemoteRunner(p.sshClient)),
		)
	} else {
		agentStartupCmds = append(agentStartupCmds,
			command.New("stop current fx agent", stopFxAgent, command.NewLocalRunner()),
			command.New("start fx agent", startFxAgent, command.NewLocalRunner()),
		)
	}
	for _, cmd := range agentStartupCmds {
		if output, err := cmd.Exec(); err != nil {
			if strings.Contains(string(output), "No such container: fx-agent") {
				// Skip stop a fx-agent error when there is not agent running
			} else {
				log.Fatalf("Provision:%s: %s, %s", cmd.Name, err, output)
				return err
			}
		}
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
