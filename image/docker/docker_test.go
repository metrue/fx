package docker

import "testing"

func TestDocker(t *testing.T) {
	cli, err := CreateClient()
	if err != nil {
		t.Fatal(err)
	}

	workdir := "./fixture"
	name := "fx-test-docker-image:latest"
	if err := cli.Build(workdir, name); err != nil {
		t.Fatal(err)
	}

	if err := cli.Push(name); err != nil {
		t.Fatal(err)
	}
}
