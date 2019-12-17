package docker

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

func TestCloudProvision(t *testing.T) {
	t.Run("fx agent started", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		n := New("127.0.0.1", "fx", "master")
		sshClient := sshMocks.NewMockClienter(ctrl)
		n.setsshClient(sshClient)
		sshClient.EXPECT().RunCommand(infra.Scripts["docker_version"].(string), ssh.CommandOptions{}).Return(nil)
		sshClient.EXPECT().RunCommand(infra.Scripts["check_fx_agent"].(string), ssh.CommandOptions{}).Return(nil)
		if err := n.Provision(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("fx agent not started", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		n := New("127.0.0.1", "fx", "master")
		sshClient := sshMocks.NewMockClienter(ctrl)
		n.setsshClient(sshClient)
		sshClient.EXPECT().RunCommand(infra.Scripts["docker_version"].(string), ssh.CommandOptions{}).Return(nil)
		sshClient.EXPECT().RunCommand(infra.Scripts["check_fx_agent"].(string), ssh.CommandOptions{}).Return(fmt.Errorf("no such container"))
		sshClient.EXPECT().RunCommand(infra.Scripts["start_fx_agent"].(string), ssh.CommandOptions{}).Return(nil)
		if err := n.Provision(); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("docker not installed and fx agent not started", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		n := New("127.0.0.1", "fx", "master")
		sshClient := sshMocks.NewMockClienter(ctrl)
		n.setsshClient(sshClient)
		sshClient.EXPECT().RunCommand(infra.Scripts["docker_version"].(string), ssh.CommandOptions{}).Return(fmt.Errorf("no such command"))
		sshClient.EXPECT().RunCommand(infra.Scripts["install_docker"].(string), ssh.CommandOptions{}).Return(nil)
		sshClient.EXPECT().RunCommand(infra.Scripts["start_dockerd"].(string), ssh.CommandOptions{}).Return(nil)
		sshClient.EXPECT().RunCommand(infra.Scripts["check_fx_agent"].(string), ssh.CommandOptions{}).Return(fmt.Errorf("no such container"))
		sshClient.EXPECT().RunCommand(infra.Scripts["start_fx_agent"].(string), ssh.CommandOptions{}).Return(nil)
		if err := n.Provision(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestCloudIsHealth(t *testing.T) {
	t.Run("agent started", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		cloud := New("127.0.0.1", "fx", "master")
		sshClient := sshMocks.NewMockClienter(ctrl)
		cloud.setsshClient(sshClient)

		sshClient.EXPECT().RunCommand(infra.Scripts["check_fx_agent"].(string), ssh.CommandOptions{}).Return(nil)
		ok, err := cloud.IsHealth()
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatalf("cloud should be healthy")
		}
	})

	t.Run("agent not started, and retart ok", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		cloud := New("127.0.0.1", "fx", "master")
		sshClient := sshMocks.NewMockClienter(ctrl)
		cloud.setsshClient(sshClient)

		sshClient.EXPECT().RunCommand(infra.Scripts["check_fx_agent"].(string), ssh.CommandOptions{}).Return(fmt.Errorf("fx agent not started"))
		sshClient.EXPECT().RunCommand(infra.Scripts["start_fx_agent"].(string), ssh.CommandOptions{}).Return(nil)
		ok, err := cloud.IsHealth()
		if err != nil {
			t.Fatal(err)
		}
		if !ok {
			t.Fatalf("cloud should be healthy")
		}
	})

	t.Run("agent not started, but restart failed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		cloud := New("127.0.0.1", "fx", "master")
		sshClient := sshMocks.NewMockClienter(ctrl)
		cloud.setsshClient(sshClient)

		sshClient.EXPECT().RunCommand(infra.Scripts["check_fx_agent"].(string), ssh.CommandOptions{}).Return(fmt.Errorf("fx agent not started"))
		sshClient.EXPECT().RunCommand(infra.Scripts["start_fx_agent"].(string), ssh.CommandOptions{}).Return(fmt.Errorf("fx agent started failed"))
		ok, err := cloud.IsHealth()
		if err == nil {
			t.Fatal("should got failed starting")
		}
		if ok {
			t.Fatalf("cloud should not be healthy")
		}
	})
}

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
