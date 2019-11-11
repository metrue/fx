package doctor

import (
	"github.com/apex/log"
	"github.com/metrue/fx/constants"
	"github.com/metrue/fx/pkg/command"
	"github.com/metrue/go-ssh-client"
)

// Doctor health checking
type Doctor struct {
	host string

	sshClient ssh.Client
}

func isLocal(host string) bool {
	if host == "" {
		return false
	}
	return host == "127.0.0.1" || host == "localhost" || host == "0.0.0.0"
}

// New a doctor
func New(host, user, password string) *Doctor {
	sshClient := ssh.New(host).
		WithUser(user).
		WithPassword(password)
	return &Doctor{
		host:      host,
		sshClient: sshClient,
	}
}

// Start diagnosis
func (d *Doctor) Start() error {
	checkDocker := "docker version"
	checkAgent := "docker inspect " + constants.AgentContainerName

	cmds := []*command.Command{}
	if !isLocal(d.host) {
		cmds = append(cmds,
			command.New("check if dockerd is running", checkDocker, command.NewRemoteRunner(d.sshClient)),
			command.New("check if fx agent is running", checkAgent, command.NewRemoteRunner(d.sshClient)),
		)
	} else {
		cmds = append(cmds,
			command.New("check if dockerd is running", checkDocker, command.NewLocalRunner()),
			command.New("check if fx agent is running", checkAgent, command.NewLocalRunner()),
		)
	}

	for _, cmd := range cmds {
		if _, err := cmd.Exec(); err != nil {
			log.Fatalf("Doctor check:%s: %s", cmd.Name, err)
		} else {
			log.Infof("Doctor check:%s: \u2713", cmd.Name)
		}
	}

	return nil
}
