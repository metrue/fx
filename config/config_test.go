package config

import (
	"os"
	"reflect"
	"testing"
)

func TestConfig(t *testing.T) {
	configPath := "/tmp/.fx"
	defer func() {
		if err := os.RemoveAll(configPath); err != nil {
			t.Fatal(err)
		}
	}()

	c := New(configPath)
	if err := c.Init(); err != nil {
		t.Fatal(err)
	}

	host, err := c.GetDefaultHost()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(host, Host{Host: "localhost"}) {
		t.Fatalf("should get %v but got %v", Host{Host: "localhost"}, host)
	}

	name := "remote-a"
	h := Host{
		Host:     "192.168.1.1",
		User:     "user-a",
		Password: "password-a",
	}
	if err := c.AddHost(name, h); err != nil {
		t.Fatal(err)
	}

	if err := c.SetDefaultHost(name, h); err != nil {
		t.Fatal(err)
	}

	host, err = c.GetDefaultHost()
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(host, h) {
		t.Fatalf("should get %v but got %v", h, host)
	}
}
