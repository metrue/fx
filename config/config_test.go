package config

import (
	"fmt"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	configPath := "./config.yml"
	defer func() {
		if err := os.RemoveAll(configPath); err != nil {
			t.Fatal(err)
		}
	}()

	os.Setenv("FX_CONFIG", configPath)
	c, err := Load()
	if err != nil {
		t.Fatal(err)
	}

	if len(c.Clouds) != 1 {
		t.Fatal("should contain default cloud")
	}

	name := "fx_cluster_1"
	if err := Use(name); err == nil {
		t.Fatal("should get no such cloud error")
	}

	cloud := Cloud{
		"type":   "k8s",
		"config": "./kubeconfig",
	}
	if err := AddCloud(name, cloud); err != nil {
		t.Fatal(err)
	}

	if err := Use(name); err != nil {
		t.Fatal(err)
	}

	body, err := View()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(body))
}
