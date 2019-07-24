package provision

import (
	"github.com/apex/log"
	ssh "github.com/metrue/go-ssh-client"
)

// Options options for provision, these information is used for ssh login and run provision script
type Options struct {
	Host     string
	User     string
	Password string
}

// Provisioner provision
type Provisioner interface {
	Start(opt Options)
}

// Provisionor provision-or
type Provisionor struct {
	sshClient ssh.Client
}

// New new provision
func New(opt Options) *Provisionor {
	sshClient := ssh.New(opt.Host).
		WithUser(opt.User).
		WithPassword(opt.Password)
	return &Provisionor{sshClient: sshClient}
}

// Start start provision progress
func (p *Provisionor) Start() error {
	if err := p.proxyDockerEngineAPI(); err != nil {
		return err
	}
	log.Info("Provision: proxy docker engine API \u2713")

	log.Info("Provision: pulling base docker images \n metrue/fx-java-base metrue/fx-julia-base metrue/fx-python-base metrue/fx-node-base metrue/fx-d-base metrue/fx-go-base ...")
	if err := p.pullBaseDockerImages(); err != nil {
		return err
	}
	log.Info("Provision: pull API \u2713")

	return nil
}

func (p *Provisionor) proxyDockerEngineAPI() error {
	// https://docs.docker.com/engine/api/v1.24/
	cmd := "docker run -d -v /var/run/docker.sock:/var/run/docker.sock -p 0.0.0.0:8866:1234 bobrik/socat TCP-LISTEN:1234,fork UNIX-CONNECT:/var/run/docker.sock"
	// cmd := "docker -v"
	if _, err := p.sshClient.RunCommand(cmd); err != nil {
		return err
	}
	return nil
}

func (p *Provisionor) pullBaseDockerImages() error {
	cmd := "docker pull metrue/fx-java-base & docker pull metrue/fx-julia-base & docker pull metrue/fx-python-base & docker pull metrue/fx-node-base & docker pull metrue/fx-d-base & docker pull metrue/fx-go-base & wait"
	if _, err := p.sshClient.RunCommand(cmd); err != nil {
		return err
	}
	return nil
}
