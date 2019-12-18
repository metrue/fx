package utils

import (
	"os"
	"path/filepath"
)

// HasDockerfile check if there is Dockerfile in dir
func HasDockerfile(dir string) bool {
	var dockerfile string
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// nolint
		if info.Mode().IsRegular() && info.Name() == "Dockerfile" {
			dockerfile = path
		}
		return nil
	}); err != nil {
		return false
	}
	if dockerfile == "" {
		return false
	}
	return true
}
