package provision

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"sync"

	"github.com/apex/log"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/metrue/fx/constants"
	ssh "github.com/metrue/go-ssh-client"
)

type containerInfo struct {
	ID         string                     `json:"Id"`
	State      dockerTypes.ContainerState `json:"State"`
	Image      string                     `json:"Image"`
	HostConfig container.HostConfig       `json:"HostConfig"`
}

const name = "docker-sock-proxy-for-fx"

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
	if !p.local {
		if err := p.enableDockerEngineAPIRemote(); err != nil {
			return err
		}
		log.Info("Provision: proxy docker engine API \u2713")

		log.Info("Provision: pulling base docker images \n metrue/fx-java-base metrue/fx-julia-base metrue/fx-python-base metrue/fx-node-base metrue/fx-d-base metrue/fx-go-base ...")
		if err := p.pullBaseDockerImagesRemote(); err != nil {
			return err
		}
		log.Info("Provision: pull base Docker images \u2713")
		return nil
	}

	if err := p.enableDockerEngineAPILocal(); err != nil {
		log.Fatalf("Provision: proxy docker engine API: %s", err)
		return err
	}
	log.Info("Provision: proxy docker engine API \u2713")

	if err := p.pullBaseDockerImagesLocal(); err != nil {
		log.Fatalf("Provision: pull base Docker images: %s", err)
		return err
	}
	log.Info("Provision: pull base Docker images \u2713")
	return nil
}

func (p *Provisionor) enableDockerEngineAPIRemote() error {
	// https://docs.docker.com/engine/api/v1.24/
	cmd := fmt.Sprintf("docker stop %s && docker rm %s && docker run -d --name=%s -v /var/run/docker.sock:/var/run/docker.sock -p 0.0.0.0:%d:1234 bobrik/socat TCP-LISTEN:1234,fork UNIX-CONNECT:/var/run/docker.sock", name, name, name, constants.AgentPort)
	// cmd := "docker -v"
	if _, err := p.sshClient.RunCommand(cmd); err != nil {
		return err
	}
	return nil
}

func (p *Provisionor) pullBaseDockerImagesRemote() error {
	cmd := "docker pull metrue/fx-java-base & docker pull metrue/fx-julia-base & docker pull metrue/fx-python-base & docker pull metrue/fx-node-base & docker pull metrue/fx-d-base & docker pull metrue/fx-go-base & wait"
	if _, err := p.sshClient.RunCommand(cmd); err != nil {
		return err
	}
	return nil
}

func (p *Provisionor) enableDockerEngineAPILocal() error {
	cmd := exec.Command(
		"docker",
		"inspect",
		name,
	)
	var infos []containerInfo
	// TODO use dockerd rest API
	// docker inspect <name>
	// it will returns:
	// []
	// Error: No such object: <name>
	// when no such container
	stdoutStderr, _ := cmd.CombinedOutput()
	if err := json.Unmarshal(stdoutStderr, &infos); err != nil {
		// no proxy container created
		return runProxy()
	}

	const containerStateCreated = "created"
	const containerStateRestarting = "restarting"
	const containerStateRunning = "running"
	const containerStatePaused = "paused"
	const containerStateExited = "exited"

	state := infos[0].State.Status
	if state == containerStateRestarting ||
		state == containerStateRunning { // FIXME should wait for it ready
		return nil
	}

	if state == containerStateCreated ||
		state == containerStatePaused ||
		state == containerStateExited {
		cmd := exec.Command(
			"docker",
			"start",
			name,
		)
		if _, err := cmd.CombinedOutput(); err != nil {
			return err
		}
	}
	return nil
}

func runProxy() error {
	// enable docker remote api
	// docker run -d -v /var/run/docker.sock:/var/run/docker.sock -p 127.0.0.1:1234:1234 bobrik/socat TCP-LISTEN:1234,fork UNIX-CONNECT:/var/run/docker.sock
	cmd := exec.Command(
		"docker",
		"run",
		"--name="+name,
		"-d",
		"-v",
		"/var/run/docker.sock:/var/run/docker.sock",
		"-p",
		fmt.Sprintf("127.0.0.1:%d:1234", constants.AgentPort),
		"bobrik/socat",
		"TCP-LISTEN:1234,fork",
		"UNIX-CONNECT:/var/run/docker.sock",
	)
	_, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	return nil
}

func (p *Provisionor) pullBaseDockerImagesLocal() error {
	var wg sync.WaitGroup
	for _, image := range constants.BaseImages {
		wg.Add(1)
		go func(img string) {
			cmd := exec.Command(
				"docker",
				"pull",
				img,
			)
			stdoutStderr, err := cmd.CombinedOutput()
			if err != nil {
				log.Fatalf("Pulling %s failed: %v", img, string(stdoutStderr))
			} else {
				log.Infof("%s Pulled: \u2713", img)
			}
			wg.Done()
		}(image)
	}
	wg.Wait()
	return nil
}
