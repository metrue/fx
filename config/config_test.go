package config

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	configPath := "./tmp/config.yml"
	defer func() {
		if err := os.RemoveAll("./tmp"); err != nil {
			t.Fatal(err)
		}
	}()
	c, err := Load(configPath)
	if err != nil {
		t.Fatal(err)
	}

	if len(c.Clouds) != 1 {
		t.Fatal("should contain default cloud")
	}

	name := "fx_cluster_1"
	if err := c.Use(name); err == nil {
		t.Fatal("should get no such cloud error")
	}

	if err := c.AddK8SCloud(name, []byte("sampe kubeconfg")); err != nil {
		t.Fatal(err)
	}

	config := map[string]string{
		"ip":   "127.0.0.1",
		"user": "use1",
	}
	configData, _ := json.Marshal(config)
	if err := c.AddDockerCloud("docker-1", configData); err != nil {
		t.Fatal(err)
	}

	if err := c.Use(name); err != nil {
		t.Fatal(err)
	}

	if c.CurrentCloud != name {
		t.Fatalf("should get %s but got %s", name, c.CurrentCloud)
	}

	conf, err := Load(configPath)
	if err != nil {
		t.Fatal(err)
	}
	if conf.CurrentCloud != name {
		t.Fatalf("should get %s but got %s", name, c.CurrentCloud)
	}

	body, err := c.View()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(body))
}
