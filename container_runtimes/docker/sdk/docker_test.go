package docker

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	dockerTypes "github.com/docker/docker/api/types"
)

func TestDocker(t *testing.T) {
	ctx := context.Background()
	cli, err := CreateClient(ctx)
	if err != nil {
		t.Fatal(err)
	}

	workdir := "../fixture"
	name := "fx-test-docker-image"
	if err := cli.BuildImage(ctx, workdir, name); err != nil {
		t.Fatal(err)
	}

	// wait a while for image to be tagged successfully after build
	time.Sleep(2 * time.Second)

	var imgInfo dockerTypes.ImageInspect
	if err := cli.InspectImage(ctx, name, &imgInfo); err != nil {
		t.Fatal(err)
	}
	found := false
	for _, t := range imgInfo.RepoTags {
		slice := strings.Split(t, ":")
		if slice[0] == name {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("should have built image with tag %s", name)
	}

	username := os.Getenv("DOCKER_USERNAME")
	password := os.Getenv("DOCKER_PASSWORD")
	if username == "" || password == "" {
		t.Skip("Skip push image test since DOCKER_USERNAME and DOCKER_PASSWORD not set in environment variable")
	}

	img, err := cli.PushImage(ctx, name)
	if err != nil {
		t.Fatal(err)
	}
	expect := username + "/" + name
	if img != expect {
		t.Fatalf("should get %s but got %s", expect, img)
	}
}
