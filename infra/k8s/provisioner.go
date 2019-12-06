package k8s

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/metrue/fx/infra"
	sshOperator "github.com/metrue/go-ssh-client"
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

// Provisioner k3s operator
type Provisioner struct {
	master MasterNode
	agents []AgentNode
}

// TODO upgrade to latest when k3s fix the tls scan issue
// https://github.com/rancher/k3s/issues/556
const version = "v0.9.1"

// New new a operator
func New(master MasterNode, agents []AgentNode) *Provisioner {
	return &Provisioner{
		master: master,
		agents: agents,
	}
}

// Provision provision k3s cluster
func (k *Provisioner) Provision() ([]byte, error) {
	if err := k.SetupMaster(); err != nil {
		return nil, err
	}
	if err := k.SetupAgent(); err != nil {
		return nil, err
	}
	return k.GetKubeConfig()
}

// HealthCheck check healthy status of host
func (k *Provisioner) HealthCheck() (bool, error) {
	// TODO
	return true, nil
}

// SetupMaster setup master node
func (k *Provisioner) SetupMaster() error {
	sshKeyFile, _ := infra.GetSSHKeyFile()
	ssh := sshOperator.New(k.master.IP).WithUser(k.master.User).WithKey(sshKeyFile)
	installCmd := fmt.Sprintf("curl -sLS https://get.k3s.io | INSTALL_K3S_EXEC='server --docker --tls-san %s' INSTALL_K3S_VERSION='%s' sh -", k.master.IP, version)
	if err := ssh.RunCommand(infra.Sudo(installCmd, k.master.User), sshOperator.CommandOptions{
		Stdout: os.Stdout,
		Stdin:  os.Stdin,
		Stderr: os.Stderr,
	}); err != nil {
		fmt.Println("setup master failed \n ===========")
		fmt.Println(err)
		fmt.Println("===========")
	}
	return nil
}

func (k *Provisioner) getToken() (string, error) {
	sshKeyFile, _ := infra.GetSSHKeyFile()
	ssh := sshOperator.New(k.master.IP).WithUser(k.master.User).WithKey(sshKeyFile)
	script := "cat /var/lib/rancher/k3s/server/node-token"
	var outPipe bytes.Buffer
	if err := ssh.RunCommand(infra.Sudo(script, k.master.User), sshOperator.CommandOptions{
		Stdout: bufio.NewWriter(&outPipe),
		Stdin:  os.Stdin,
		Stderr: os.Stderr,
	}); err != nil {
		return "", err
	}
	return outPipe.String(), nil
}

// SetupAgent set agent node
func (k *Provisioner) SetupAgent() error {
	sshKeyFile, _ := infra.GetSSHKeyFile()
	tok, err := k.getToken()
	if err != nil {
		return err
	}
	const k3sExtraArgs = "--docker"
	joinCmd := fmt.Sprintf("curl -fL https://get.k3s.io/ | K3S_URL='https://%s:6443' K3S_TOKEN='%s' INSTALL_K3S_VERSION='%s' sh -s - %s", k.master.IP, tok, version, k3sExtraArgs)
	for _, agent := range k.agents {
		ssh := sshOperator.New(agent.IP).WithUser(agent.User).WithKey(sshKeyFile)
		if err := ssh.RunCommand(joinCmd, sshOperator.CommandOptions{
			Stdout: os.Stdout,
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
		}); err != nil {
			fmt.Println("setup agent failed \n================")
			fmt.Println(err)
			fmt.Println("================")
			return err
		}
	}

	return nil
}

// GetKubeConfig get kubeconfig of k3s cluster
func (k *Provisioner) GetKubeConfig() ([]byte, error) {
	sshKeyFile, _ := infra.GetSSHKeyFile()
	var config []byte
	getConfigCmd := "cat /etc/rancher/k3s/k3s.yaml\n"
	ssh := sshOperator.New(k.master.IP).WithUser(k.master.User).WithKey(sshKeyFile)
	var outPipe bytes.Buffer
	if err := ssh.RunCommand(infra.Sudo(getConfigCmd, k.master.User), sshOperator.CommandOptions{
		Stdout: bufio.NewWriter(&outPipe),
		Stdin:  os.Stdin,
		Stderr: os.Stderr,
	}); err != nil {
		fmt.Println("setup agent failed \n================")
		fmt.Println("================")
		fmt.Println(err)
		return config, err
	}
	return rewriteKubeconfig(outPipe.String(), k.master.IP, "default"), nil
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

var _ infra.Provisioner = &Provisioner{}
