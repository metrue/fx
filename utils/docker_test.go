package utils

import "testing"

func TestHasDockerfile(t *testing.T) {
	dir := "tmp"
	_ = EnsureDir(dir)

	if HasDockerfile(dir) {
		t.Fatalf("should get false but got true")
	}
}
