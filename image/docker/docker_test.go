package docker

import "testing"

func TestDocker(t *testing.T) {
	cli, err := CreateClient()
	if err != nil {
		t.Fatal(err)
	}

	workdir := "./fixture"
	name := "fx-test-docker-iamge-builder"
	if err := cli.Build(workdir, name); err != nil {
		t.Fatal(err)
	}
}
