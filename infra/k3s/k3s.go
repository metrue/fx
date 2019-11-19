package k3s

import (
	"fmt"
	"strings"

	sshOperator "github.com/metrue/go-ssh-client"
	"github.com/mitchellh/go-homedir"
)

// MasterNode master node instance
type MasterNode struct {
	IP   string
	User string
}

// AgentNode agent node instance
type AgentNode struct {
	IP   string
	User string
}

// K3S k3s operator
type K3S struct {
	master MasterNode
	agents []AgentNode
}

// TODO upgrade to latest when k3s fix the tls scan issue
// https://github.com/rancher/k3s/issues/556
const version = "v0.9.1"

// New new a operator
func New(master MasterNode, agents []AgentNode) *K3S {
	return &K3S{
		master: master,
		agents: agents,
	}
}

// SetupMaster setup master node
func (k *K3S) SetupMaster() error {
	publicKey, err := homedir.Expand("~/.ssh/id_rsa")
	if err != nil {
		return err
	}
	ssh := sshOperator.New(k.master.IP).WithUser(k.master.User).WithKey(publicKey)
	installCmd := fmt.Sprintf("curl -sLS https://get.k3s.io | INSTALL_K3S_EXEC='server --tls-san %s' INSTALL_K3S_VERSION='%s' sh -", k.master.IP, version)
	stdout, stderr, err := ssh.RunCommand(installCmd)
	if err != nil {
		fmt.Println("setup master failed \n ===========")
		fmt.Println("failed: ", string(stderr))
		fmt.Println("output: ", string(stdout))
		fmt.Println("===========")
	}
	return err
}

func (k *K3S) getToken() (string, error) {
	publicKey, err := homedir.Expand("~/.ssh/id_rsa")
	if err != nil {
		return "", err
	}
	ssh := sshOperator.New(k.master.IP).WithUser(k.master.User).WithKey(publicKey)
	script := "cat /var/lib/rancher/k3s/server/node-token"
	stdout, _, err := ssh.RunCommand(script)
	if err != nil {
		return "", err
	}
	return string(stdout), nil
}

// SetupAgent set agent node
func (k *K3S) SetupAgent() error {
	publicKey, err := homedir.Expand("~/.ssh/id_rsa")
	if err != nil {
		return err
	}
	tok, err := k.getToken()
	if err != nil {
		return err
	}
	const k3sExtraArgs = ""
	joinCmd := fmt.Sprintf("curl -sfL https://get.k3s.io/ | K3S_URL='https://%s:6443' K3S_TOKEN='%s' INSTALL_K3S_VERSION='%s' sh -s - %s", k.master.IP, tok, version, k3sExtraArgs)
	for _, agent := range k.agents {
		ssh := sshOperator.New(agent.IP).WithUser(agent.User).WithKey(publicKey)
		stdout, stderr, err := ssh.RunCommand(joinCmd)
		if err != nil {
			fmt.Println("setup agent failed \n================")
			fmt.Println("failed: ", string(stderr))
			fmt.Println("output: ", string(stdout))
			fmt.Println("================")
			return err
		}
		return nil
	}

	return nil
}

// GetKubeConfig get kubeconfig of k3s cluster
func (k *K3S) GetKubeConfig() ([]byte, error) {
	publicKey, err := homedir.Expand("~/.ssh/id_rsa")
	var config []byte
	if err != nil {
		return config, err
	}
	getConfigCmd := "cat /etc/rancher/k3s/k3s.yaml\n"
	ssh := sshOperator.New(k.master.IP).WithUser(k.master.User).WithKey(publicKey)
	stdout, stderr, err := ssh.RunCommand(getConfigCmd)
	if err != nil {
		fmt.Println("setup agent failed \n================")
		fmt.Println("failed: ", string(stderr))
		fmt.Println("output: ", string(stdout))
		fmt.Println("================")
		return config, err
	}
	return rewriteKubeconfig(string(stdout), k.master.IP, "default"), nil
}

func rewriteKubeconfig(kubeconfig string, ip string, context string) []byte {
	if context == "" {
		context = "default"
	}

	kubeconfigReplacer := strings.NewReplacer(
		"127.0.0.1", ip,
		"localhost", ip,
		"default", context,
	)

	return []byte(kubeconfigReplacer.Replace(kubeconfig))
}
