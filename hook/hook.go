package hook

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/metrue/fx/utils"
)

// Hooker defines hook interface
type Hooker interface {
	Run() error
}

// Hook to run
type Hook struct {
	name   string
	script string
}

// New a hook
func New(name string, script string, workdir string) *Hook {
	return &Hook{
		name:   name,
		script: script,
	}
}

// Run execute a hook
func (h *Hook) Run(workdir string) error {
	var script string
	if !utils.IsRegularFile(h.script) {
		hookScript, err := ioutil.TempFile(os.TempDir(), "fx-hook-script-")
		if err != nil {
			return err
		}
		defer os.Remove(hookScript.Name())

		content := []byte(h.script)
		if _, err = hookScript.Write(content); err != nil {
			return err
		}
		if err := hookScript.Close(); err != nil {
			return err
		}
		script = hookScript.Name()
	} else {
		absScript, err := filepath.Abs(h.script)
		if err != nil {
			return err
		}
		script = absScript
	}

	cmd := exec.Command("/bin/sh", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if workdir != "" {
		cmd.Dir = workdir
	}

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// Name hook name
func (h *Hook) Name() string {
	return h.name
}
