package infra

import (
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
)

// GetSSHKeyFile get ssh private key file
func GetSSHKeyFile() (string, error) {
	path := os.Getenv("SSH_KEY_FILE")
	if path != "" {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", err
		}
		return absPath, nil
	}

	key, err := homedir.Expand("~/.ssh/id_rsa")
	if err != nil {
		return "", err
	}
	return key, nil
}

// GetSSHPort get ssh port
func GetSSHPort() string {
	port := os.Getenv("SSH_PORT")
	if port != "" {
		return port
	}
	return "22"
}
