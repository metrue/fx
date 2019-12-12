package k8s

import (
	"fmt"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/metrue/fx/infra"
	"github.com/metrue/go-ssh-client"
	sshMocks "github.com/metrue/go-ssh-client/mocks"
	"github.com/mitchellh/go-homedir"
)

func TestGetSSHKeyFile(t *testing.T) {
	t.Run("defaut", func(t *testing.T) {
		defau, err := sshkey()
		if err != nil {
			t.Fatal(err)
		}

		defaultPath, _ := homedir.Expand("~/.ssh/id_rsa")
		if defau != defaultPath {
			t.Fatalf("should get %s but got %s", defaultPath, defau)
		}
	})

	t.Run("override from env", func(t *testing.T) {
		os.Setenv("SSH_KEY_FILE", "/tmp/id_rsa")
		keyFile, err := sshkey()
		if err != nil {
			t.Fatal(err)
		}
		if keyFile != "/tmp/id_rsa" {
			t.Fatalf("should get %s but got %s", "/tmp/id_rsa", keyFile)
		}
	})
}

func TestGetSSHPort(t *testing.T) {
	t.Run("defaut", func(t *testing.T) {
		defau := sshport()
		if defau != "22" {
			t.Fatalf("should get %s but got %s", "22", defau)
		}
	})

	t.Run("override from env", func(t *testing.T) {
		os.Setenv("SSH_PORT", "2222")
		defau := sshport()
		if defau != "2222" {
			t.Fatalf("should get %s but got %s", "2222", defau)
		}
	})
}

func TestNode(t *testing.T) {
	t.Run("master node already has docker and k3s server", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		n, err := CreateNode("127.0.0.1", "fx", NodeTypeMaster, "master")
		if err != nil {
			t.Fatal(err)
		}

		if n.sshClient == nil {
			t.Fatal("ssh client should not be nil")
		}

		sshClient := sshMocks.NewMockClienter(ctrl)
		n.setsshClient(sshClient)
		sshClient.EXPECT().RunCommand(infra.Scripts["docker_version"].(string), ssh.CommandOptions{
			Stdout: os.Stdout,
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
		}).Return(nil)
		sshClient.EXPECT().RunCommand(infra.Scripts["check_k3s_server"].(string), ssh.CommandOptions{
			Stdout: os.Stdout,
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
		}).Return(nil)
		if err := n.Provision(map[string]string{}); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("master node no docker and k3s server", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		n, err := CreateNode("127.0.0.1", "fx", NodeTypeMaster, "master")
		if err != nil {
			t.Fatal(err)
		}

		if n.sshClient == nil {
			t.Fatal("ssh client should not be nil")
		}

		sshClient := sshMocks.NewMockClienter(ctrl)
		n.setsshClient(sshClient)
		sshClient.EXPECT().RunCommand(infra.Scripts["docker_version"].(string), ssh.CommandOptions{
			Stdout: os.Stdout,
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
		}).Return(fmt.Errorf("no such command"))
		sshClient.EXPECT().RunCommand(infra.Scripts["install_docker"].(string), ssh.CommandOptions{
			Stdout: os.Stdout,
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
		}).Return(nil)
		sshClient.EXPECT().RunCommand(infra.Scripts["start_dockerd"].(string), ssh.CommandOptions{
			Stdout: os.Stdout,
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
		}).Return(nil)
		sshClient.EXPECT().RunCommand(infra.Scripts["check_k3s_server"].(string), ssh.CommandOptions{
			Stdout: os.Stdout,
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
		}).Return(fmt.Errorf("no such progress"))

		cmd := infra.Scripts["setup_k3s_master"].(func(ip string) string)(n.IP)
		sshClient.EXPECT().RunCommand(cmd, ssh.CommandOptions{
			Stdout: os.Stdout,
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
		}).Return(nil)
		if err := n.Provision(map[string]string{}); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("agent node already has docker and k3s agent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		n, err := CreateNode("127.0.0.1", "fx", NodeTypeAgent, "agent")
		if err != nil {
			t.Fatal(err)
		}

		if n.sshClient == nil {
			t.Fatal("ssh client should not be nil")
		}

		sshClient := sshMocks.NewMockClienter(ctrl)
		n.setsshClient(sshClient)
		sshClient.EXPECT().RunCommand(infra.Scripts["docker_version"].(string), ssh.CommandOptions{
			Stdout: os.Stdout,
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
		}).Return(nil)
		sshClient.EXPECT().RunCommand(infra.Scripts["check_k3s_agent"].(string), ssh.CommandOptions{
			Stdout: os.Stdout,
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
		}).Return(nil)
		if err := n.Provision(map[string]string{}); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("agent node no docker and k3s agent", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		n, err := CreateNode("127.0.0.1", "fx", NodeTypeAgent, "agent")
		if err != nil {
			t.Fatal(err)
		}

		if n.sshClient == nil {
			t.Fatal("ssh client should not be nil")
		}

		sshClient := sshMocks.NewMockClienter(ctrl)
		n.setsshClient(sshClient)
		sshClient.EXPECT().RunCommand(infra.Scripts["docker_version"].(string), ssh.CommandOptions{
			Stdout: os.Stdout,
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
		}).Return(fmt.Errorf("no such command"))
		sshClient.EXPECT().RunCommand(infra.Scripts["install_docker"].(string), ssh.CommandOptions{
			Stdout: os.Stdout,
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
		}).Return(nil)
		sshClient.EXPECT().RunCommand(infra.Scripts["start_dockerd"].(string), ssh.CommandOptions{
			Stdout: os.Stdout,
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
		}).Return(nil)
		sshClient.EXPECT().RunCommand(infra.Scripts["check_k3s_agent"].(string), ssh.CommandOptions{
			Stdout: os.Stdout,
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
		}).Return(fmt.Errorf("no such progress"))

		url := "url-1"
		token := "token-1"
		cmd := infra.Scripts["setup_k3s_agent"].(func(url string, ip string) string)(url, token)
		sshClient.EXPECT().RunCommand(cmd, ssh.CommandOptions{
			Stdout: os.Stdout,
			Stdin:  os.Stdin,
			Stderr: os.Stderr,
		}).Return(nil)
		if err := n.Provision(map[string]string{"url": url, "token": token}); err != nil {
			t.Fatal(err)
		}
	})
}
