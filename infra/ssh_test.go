package infra

import (
	"os"
	"testing"

	"github.com/mitchellh/go-homedir"
)

func TestGetSSHKeyFile(t *testing.T) {
	t.Run("defaut", func(t *testing.T) {
		defau, err := GetSSHKeyFile()
		if err != nil {
			t.Fatal(err)
		}

		defaultPath, _ := homedir.Expand("~/.ssh/id_rsa")
		if defau != defaultPath {
			t.Fatalf("should get %s but got %s", defaultPath, defau)
		}
	})

	t.Run("override from env", func(t *testing.T) {
		os.Setenv("SSH_KEY_FILE", "/tmp/id_rsa")
		keyFile, err := GetSSHKeyFile()
		if err != nil {
			t.Fatal(err)
		}
		if keyFile != "/tmp/id_rsa" {
			t.Fatalf("should get %s but got %s", "/tmp/id_rsa", keyFile)
		}
	})
}
