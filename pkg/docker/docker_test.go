package docker_test

import (
	"os"
	"os/exec"
	"testing"

	. "github.com/metrue/fx/pkg/docker"
	"github.com/stretchr/testify/assert"
)

func TestInfo(t *testing.T) {
	ok := IsRunning()
	assert.Equal(t, true, ok)
}

func checkDocker() {
	if err := exec.Command("docker", "info").Run(); err != nil {
		panic(os.Stderr)
	}
}

func TestMain(m *testing.M) {
	checkDocker()
	m.Run()
}
