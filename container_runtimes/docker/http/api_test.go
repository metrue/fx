package api

import (
	"context"
	"os"
	"testing"

	"github.com/docker/docker/api/types"
)

func TestDockerHTTP(t *testing.T) {
	host := os.Getenv("DOCKER_ENGINE_HOST")
	port := os.Getenv("DOCKER_ENGINE_PORT")
	if host == "" ||
		port == "" {
		t.Skip("DOCKER_ENGINE_HOST and DOCKER_ENGINE_PORT required")
	}

	api, err := Create(host, port)
	if err != nil {
		t.Fatal(err)
	}
	name := "fx-agent"
	var container types.ContainerJSON
	if err := api.InspectContainer(context.Background(), name, &container); err != nil {
		t.Fatal(err)
	}
	if container.Name != "/"+name {
		t.Fatalf("should get %s but got %s", name, container.Name)
	}
}
