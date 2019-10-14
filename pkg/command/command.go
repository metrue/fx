package command

import (
	"os/exec"
	"strings"

	"github.com/metrue/go-ssh-client"
)

// Runner to exec command
type Runner interface {
	Run(script string) ([]byte, error)
}

// RemoteRunner remote runner
type RemoteRunner struct {
	sshClient ssh.Client
}

// NewRemoteRunner create a runner
func NewRemoteRunner(sshClient ssh.Client) *RemoteRunner {
	return &RemoteRunner{sshClient: sshClient}
}

// Run script on remote host
func (r *RemoteRunner) Run(script string) ([]byte, error) {
	stdout, stderr, err := r.sshClient.RunCommand(script)
	output := string(stdout) + string(stderr)
	return []byte(output), err
}

// LocalRunner local runner
type LocalRunner struct{}

// NewLocalRunner create a local runner
func NewLocalRunner() *LocalRunner {
	return &LocalRunner{}
}

// Run script on local host
func (l *LocalRunner) Run(script string) ([]byte, error) {
	params := strings.Split(script, " ")
	var cmd *exec.Cmd
	if len(params) > 1 {
		// nolint: gosec
		cmd = exec.Command(params[0], params[1:]...)
	} else {
		// nolint: gosec
		cmd = exec.Command(params[0])
	}
	return cmd.CombinedOutput()
}

// Commander command interface
type Commander interface {
	Exec() ([]byte, error)
}

// Command a command
type Command struct {
	Name   string
	script string
	runner Runner
}

// New create a command
func New(name string, script string, runner Runner) *Command {
	return &Command{
		Name:   name,
		script: script,
		runner: runner,
	}
}

// Exec run command
func (c *Command) Exec() ([]byte, error) {
	return c.runner.Run(c.script)
}
