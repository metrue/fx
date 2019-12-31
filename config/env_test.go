package config

import (
	"os"
	"testing"
)

var _ = func() (_ struct{}) {
	os.Setenv("DISABLE_CONTAINER_AUTOREMOVE", "true")
	return
}()

func TestEnvLoad(t *testing.T) {
	if !DisableContainerAutoremove {
		t.Fatalf("should be true after set")
	}
}
