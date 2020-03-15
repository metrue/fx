package packer

import (
	"testing"
)

func TestTreeAndUnTree(t *testing.T) {
	_, err := PackIntoK8SConfigMapFile("./fixture/p1")
	if err != nil {
		t.Fatal(err)
	}
}
