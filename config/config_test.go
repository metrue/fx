package config

import (
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	configPath := "/tmp/.fx"
	defer func() {
		if err := os.RemoveAll(configPath); err != nil {
			t.Fatal(err)
		}
	}()

	if err := Init(configPath); err != nil {
		t.Fatal(err)
	}
	host := "localhost"
	if err := SetHost(host); err != nil {
		t.Fatal(err)
	}

	if h := GetHost(); h != host {
		t.Fatalf("should get %s but got %s", host, h)
	}
}
