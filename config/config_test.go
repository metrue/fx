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

	hosts, err := c.ListMachines()
	if err != nil {
		t.Fatal(err)
	}

	if len(hosts) != 1 {
		t.Fatalf("should have localhost as default machine")
	}

	host := hosts["localhost"]
	if !reflect.DeepEqual(host, Host{Host: "localhost", Enabled: true}) {
		t.Fatalf("should get %v but got %v", Host{Host: "localhost"}, host)
	}

	name := "remote-a"
	h := Host{
		Host:     "192.168.1.1",
		User:     "user-a",
		Password: "password-a",
		Enabled:  false,
	}
	if err := c.AddMachine(name, h); err != nil {
		t.Fatal(err)
	}

	hosts, err = c.ListMachines()
	if err != nil {
		t.Fatal(err)
	}
	if len(hosts) != 2 {
		t.Fatalf("should have %d machines now, but got %d", 2, len(hosts))
	}

	lst, err := c.ListActiveMachines()
	if err != nil {
		t.Fatal(err)
	}

	if len(lst) != 1 {
		t.Fatalf("should only have %d machine enabled, but got %d", 1, len(lst))
	}

	if err := c.EnableMachine(name); err != nil {
		t.Fatal(err)
	}

	lst, err = c.ListActiveMachines()
	if err != nil {
		t.Fatal(err)
	}

	if len(lst) != 2 {
		t.Fatalf("should only have %d machine enabled, but got %d", 2, len(lst))
	}

	h.Enabled = true
	if !reflect.DeepEqual(lst[name], h) {
		t.Fatalf("should get %v but got %v", h, lst[name])
	}

	if lst[name].Provisioned != false {
		t.Fatalf("should get %v but got %v", false, lst[name].Provisioned)
	}

	if err := c.UpdateProvisionedStatus(name, true); err != nil {
		t.Fatal(err)
	}

	updatedHost, err := c.GetMachine(name)
	if err != nil {
		t.Fatal(err)
	}

	if updatedHost.Provisioned != true {
		t.Fatalf("should get %v but got %v", true, updatedHost.Provisioned)
	}
}
