package docker

import (
	"os"
	"testing"
)

func TestDocker(t *testing.T) {
	cli, err := CreateClient()
	if err != nil {
		t.Fatal(err)
	}

	workdir := "./fixture"
	name := "fx-test-docker-image"
	if err := cli.Build(workdir, name); err != nil {
		t.Fatal(err)
	}
	username := os.Getenv("DOCKER_USERNAME")
	password := os.Getenv("DOCKER_PASSWORD")
	if username == "" || password == "" {
		t.Skip("Skip push image test since DOCKER_USERNAME and DOCKER_PASSWORD not set in enviroment variable")
	}

	img, err := cli.Push(name)
	if err != nil {
		t.Fatal(err)
	}
	expect := username + "/" + name
	if img != expect {
		t.Fatalf("should get %s but got %s", expect, img)
	}
}
