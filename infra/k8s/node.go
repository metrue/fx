package k8s

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/metrue/fx/infra"
	"github.com/metrue/go-ssh-client"
	"github.com/mitchellh/go-homedir"
)

const NodeTypeMaster = "k3s_master"
const NodeTypeAgent = "k3s_agent"
const NodeTypeDocker = "docker_agent"

// Noder node interface
type Noder interface {
	Provision(meta map[string]string) error
	GetConfig() (string, error)
	GetType() string
	GetName() string
	GetUser() string
	GetToken() (string, error)
	GetIP() string
	Dump() map[string]string
}

// Node define a node
type Node struct {
	IP   string `json:"ip"`
	User string `json:"user"`

	Type string `json:"type"`
	Name string `json:"name"`

	sshClient ssh.Clienter
}

// CreateNode create a node
func CreateNode(ip string, user string, typ string, name string) (*Node, error) {
	key, err := sshkey()
	if err != nil {
		return nil, err
	}
	port := sshport()
	sshClient := ssh.New(ip).WithUser(user).WithKey(key).WithPort(port)

	return &Node{
		IP:   ip,
		User: user,
		Type: typ,
		Name: name,

		sshClient: sshClient,
	}, nil
}

func (n *Node) runCmd(script string) error {
	return n.sshClient.RunCommand(script, ssh.CommandOptions{
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
		Stderr: os.Stderr,
	})
}

// Provision provision node
func (n *Node) Provision(meta map[string]string) error {
	if err := n.runCmd(infra.Scripts["docker_version"].(string)); err != nil {
		if err := n.runCmd(infra.Scripts["install_docker"].(string)); err != nil {
			return err
		}

		if err := n.runCmd(infra.Scripts["start_dockerd"].(string)); err != nil {
			return err
		}
	}

	if n.Type == NodeTypeMaster {
		if err := n.runCmd(infra.Scripts["check_k3s_server"].(string)); err != nil {
			cmd := infra.Scripts["setup_k3s_master"].(func(ip string) string)(n.IP)
			if err := n.runCmd(cmd); err != nil {
				return err
			}
		}
	} else if n.Type == NodeTypeAgent {
		if err := n.runCmd(infra.Scripts["check_k3s_agent"].(string)); err != nil {
			cmd := infra.Scripts["setup_k3s_agent"].(func(url string, tok string) string)(meta["url"], meta["token"])
			if err := n.runCmd(cmd); err != nil {
				return err
			}
		}
	}

	return nil
}

// GetToken get token from master node
func (n *Node) GetToken() (string, error) {
	if n.Type != NodeTypeMaster {
		return "", fmt.Errorf("could not get token from a non-master node")
	}
	var outPipe bytes.Buffer
	if err := n.sshClient.RunCommand(infra.Scripts["get_k3s_token"].(string), ssh.CommandOptions{Stdout: bufio.NewWriter(&outPipe)}); err != nil {
		return "", err
	}
	return outPipe.String(), nil
}

// State get node state
func (n *Node) State() {}

// Dump node information to json
func (n *Node) Dump() map[string]string {
	return map[string]string{
		"ip":   n.IP,
		"name": n.Name,
		"user": n.User,
		"type": n.Type,
	}
}

// GetType get node type
func (n *Node) GetType() string {
	return n.Type
}

// GetName get node type
func (n *Node) GetName() string {
	return n.Name
}

// GetIP get node type
func (n *Node) GetIP() string {
	return n.IP
}

// GetUser get user
func (n *Node) GetUser() string {
	return n.User
}

// GetConfig get config
func (n *Node) GetConfig() (string, error) {
	if n.Type == NodeTypeMaster {
		var outPipe bytes.Buffer
		if err := n.sshClient.RunCommand(infra.Scripts["get_k3s_kubeconfig"].(string), ssh.CommandOptions{
			Stdout: bufio.NewWriter(&outPipe),
		}); err != nil {
			return "", err
		}
		return string(rewriteKubeconfig(outPipe.String(), n.IP, "default")), nil
	} else if n.Type == NodeTypeDocker {
		data, err := json.Marshal(n.Dump())
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	return "", fmt.Errorf("no config for node type of %s", n.Type)
}

// NOTE only using for unit testing
func (n *Node) setsshClient(client ssh.Clienter) {
	n.sshClient = client
}

// NOTE the reason putting sshkey() and sshport here inside node.go is because
// ssh key and ssh port is related to node it self, we may extend this in future
func sshkey() (string, error) {
	path := os.Getenv("SSH_KEY_FILE")
	if path != "" {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", err
		}
		return absPath, nil
	}

	key, err := homedir.Expand("~/.ssh/id_rsa")
	if err != nil {
		return "", err
	}
	return key, nil
}

func sshport() string {
	port := os.Getenv("SSH_PORT")
	if port != "" {
		return port
	}
	return "22"
}

func rewriteKubeconfig(kubeconfig string, ip string, context string) []byte {
	if context == "" {
		// nolint
		context = "default"
	}

	kubeconfigReplacer := strings.NewReplacer(
		"127.0.0.1", ip,
		"localhost", ip,
		"default", context,
	)

	return []byte(kubeconfigReplacer.Replace(kubeconfig))
}

var (
	_ Noder = &Node{}
)
