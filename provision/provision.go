package provision

import (
	"fmt"
	"sync"

	"github.com/apex/log"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/metrue/fx/config"
	"github.com/metrue/fx/constants"
	"github.com/metrue/fx/pkg/command"
	ssh "github.com/metrue/go-ssh-client"
)

type containerInfo struct {
	ID         string                     `json:"Id"`
	State      dockerTypes.ContainerState `json:"State"`
	Image      string                     `json:"Image"`
	HostConfig container.HostConfig       `json:"HostConfig"`
}

// Options options for provision, these information is used for ssh login and run provision script
type Options struct {
	Host     string
	User     string
	Password string
}

// Provisioner provision
type Provisioner interface {
	Start()
}

// Provisionor provision-or
type Provisionor struct {
	sshClient ssh.Client
	local     bool
}

// New new provision
func New(opt Options) *Provisionor {
	local := opt.Host == "127.0.0.1" || opt.Host == "localhost"

	sshClient := ssh.New(opt.Host).
		WithUser(opt.User).
		WithPassword(opt.Password)
	return &Provisionor{
		sshClient: sshClient,
		local:     local,
	}
}

// Start start provision progress
func (p *Provisionor) Start() error {
	proxyScript := fmt.Sprintf("docker run -d --name=%s --rm -v /var/run/docker.sock:/var/run/docker.sock -p 0.0.0.0:%s:1234 bobrik/socat TCP-LISTEN:1234,fork UNIX-CONNECT:/var/run/docker.sock", constants.AgentContainerName, constants.AgentPort)
	scripts := map[string]string{
		"proxy Docker engine API":       proxyScript,
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
			if config.IsRemote() {
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
