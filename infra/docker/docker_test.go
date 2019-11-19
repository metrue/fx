package docker

import (
	"os"
	"testing"
)

func TestDocker(t *testing.T) {
	if os.Getenv("DOCKER_HOST") == "" ||
		os.Getenv("DOCKER_USER") == "" {
		t.Skip("skip test since DOCKER_HOST and DOCKER_USER not ready")
	}
	d := New(os.Getenv("DOCKER_HOST"), os.Getenv("DOCKER_USER"))
	if err := d.Install(); err != nil {
		t.Fatal(err)
	}
	if err := d.StartDockerd(); err != nil {
		t.Fatal(err)
	}
	if err := d.StartFxAgent(); err != nil {
		t.Fatal(err)
	}
}
