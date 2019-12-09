package config

import (
	"os"
	"testing"
)

func TestContainer(t *testing.T) {
	configPath := "./tmp/container.yml"
	defer func() {
		if err := os.RemoveAll("./tmp/container.yml"); err != nil {
			t.Fatal(err)
		}
	}()
	c, err := CreateContainer(configPath)
	if err != nil {
		t.Fatal(err)
	}

	if err := c.set("", ""); err == nil {
		t.Fatalf("should get error when key is empty")
	}

	if c.get("1") != nil {
		t.Fatalf("should get %v but got %v", nil, c.get("key"))
	}

	// create
	if err := c.set("1", "1"); err != nil {
		t.Fatal(err)
	}

	// read
	if c.get("1").(string) != "1" {
		t.Fatalf("should get %s but got %s", "val-1", c.get("key"))
	}

	// invaliad set
	if err := c.set("1.1", "1.1"); err != nil {
		t.Fatal(err)
	}
	if c.get("1.1").(string) != "1.1" {
		t.Fatalf("should get 1.1 but got %s", c.get("1.1"))
	}

	// update
	if err := c.set("1", "11"); err != nil {
		t.Fatal(err)
	}
	if c.get("1").(string) != "11" {
		t.Fatalf("should get 11 but got %s", c.get("1").(string))
	}

	// nested set
	if err := c.set("2.2.2.2", "2222"); err == nil {
		t.Fatalf("should throw error since 2.2.2 not ready yet")
	}

	if err := c.set("2", map[string]interface{}{
		"2": map[string]interface{}{
			"2": "2",
		},
	}); err != nil {
		t.Fatal(err)
	}

	if c.get("2.2.2").(string) != "2" {
		t.Fatalf("should get 2 but got %s", c.get("2.2.2"))
	}
	if err := c.set("2.2.2.2", "2222"); err != nil {
		t.Fatal(err)
	}
	if c.get("2.2.2.2").(string) != "2222" {
		t.Fatalf("should get 2222 but got %s", c.get("2.2.2.2"))
	}

	if err := c.set("2.2.2.1", "1111"); err != nil {
		t.Fatal(err)
	}

	if c.get("2.2.2.1").(string) != "1111" {
		t.Fatalf("should get 1111 but got %s", c.get("2.2.2.1"))
	}
}
