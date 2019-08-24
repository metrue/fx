package utils

import "testing"

func TestDockerVersion(t *testing.T) {
	host := "127.0.0.1"
	port := "8866"
	version, err := DockerVersion(host, port)
	if err != nil {
		t.Fatal(err)
	}
	if version == "" {
		t.Fatal("should version empty")
	}
}
