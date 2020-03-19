package darwin

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/metrue/go-ssh-client"
	sshMocks "github.com/metrue/go-ssh-client/mocks"
)

func TestDriverProvision(t *testing.T) {
	t.Run("SSHConnectError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sshClient := sshMocks.NewMockClienter(ctrl)
		n := &Docker{sshClient: sshClient}
		err := errors.New("could not connect to host")
		sshClient.EXPECT().Connectable(sshConnectionTimeout).Return(false, err).AnyTimes()
		if err := n.Provision(context.Background(), true); err == nil {
			t.Fatalf("should get error when SSH connection not ok")
		}
	})

	t.Run("SSHConnectionNotOK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sshClient := sshMocks.NewMockClienter(ctrl)
		n := New(sshClient)
		sshClient.EXPECT().Connectable(sshConnectionTimeout).Return(false, nil).AnyTimes()
		if err := n.Provision(context.Background(), true); err == nil {
			t.Fatalf("should get error when SSH connection not ok")
		}
	})

	t.Run("DockerAndFxAgentOK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sshClient := sshMocks.NewMockClienter(ctrl)
		n := New(sshClient)
		sshClient.EXPECT().Connectable(sshConnectionTimeout).Return(true, nil).AnyTimes()
		sshClient.EXPECT().RunCommand(scripts["docker_version"], ssh.CommandOptions{Timeout: sshConnectionTimeout}).Return(nil)
		sshClient.EXPECT().RunCommand(scripts["check_fx_agent"], ssh.CommandOptions{Timeout: sshConnectionTimeout}).Return(nil)
		if err := n.Provision(context.Background(), true); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("DockerNotReady", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sshClient := sshMocks.NewMockClienter(ctrl)
		n := New(sshClient)
		sshClient.EXPECT().Connectable(sshConnectionTimeout).Return(true, nil).AnyTimes()
		err := errors.New("docker command not found")
		sshClient.EXPECT().RunCommand(scripts["docker_version"], ssh.CommandOptions{Timeout: sshConnectionTimeout}).Return(err)
		sshClient.EXPECT().RunCommand(scripts["has_docker"], ssh.CommandOptions{Timeout: sshConnectionTimeout}).Return(err)
		if err := n.Provision(context.Background(), true); err == nil {
			t.Fatal("should tell user to install docker first")
		}
	})

	t.Run("FxAgentNotReady", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sshClient := sshMocks.NewMockClienter(ctrl)
		n := New(sshClient)
		sshClient.EXPECT().Connectable(sshConnectionTimeout).Return(true, nil).AnyTimes()
		err := errors.New("fx agent not found")
		sshClient.EXPECT().RunCommand(scripts["docker_version"], ssh.CommandOptions{Timeout: sshConnectionTimeout}).Return(nil)
		sshClient.EXPECT().RunCommand(scripts["check_fx_agent"], ssh.CommandOptions{Timeout: sshConnectionTimeout}).Return(err)
		sshClient.EXPECT().RunCommand(scripts["start_fx_agent"], ssh.CommandOptions{Timeout: sshConnectionTimeout}).Return(nil)
		if err := n.Provision(context.Background(), true); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("DockerAndFxAgentReady", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		sshClient := sshMocks.NewMockClienter(ctrl)
		n := New(sshClient)

		sshClient.EXPECT().Connectable(sshConnectionTimeout).Return(true, nil).AnyTimes()
		sshClient.EXPECT().RunCommand(scripts["docker_version"], ssh.CommandOptions{Timeout: sshConnectionTimeout}).Return(nil)
		sshClient.EXPECT().RunCommand(scripts["check_fx_agent"], ssh.CommandOptions{Timeout: sshConnectionTimeout}).Return(nil)
		if err := n.Provision(context.Background(), true); err != nil {
			t.Fatal(err)
		}
	})
}

func TestRunCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	sshClient := sshMocks.NewMockClienter(ctrl)
	n := &Docker{
		sshClient: sshClient,
	}
	script := "script"
	option := ssh.CommandOptions{
		Timeout: sshConnectionTimeout,
	}
	sshClient.EXPECT().Connectable(sshConnectionTimeout).Return(true, nil)
	sshClient.EXPECT().RunCommand(script, option).Return(nil)
	if err := n.runCmd(script, true, option); err != nil {
		t.Fatal(err)
	}
}
