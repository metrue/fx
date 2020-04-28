package packer

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/metrue/fx/utils"
)

// PackIntoK8SConfigMapFile pack function a K8S config map file
func PackIntoK8SConfigMapFile(dir string) (string, error) {
	tree := map[string]string{}
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relpath := strings.Replace(path, dir, "", 1)
			body, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			tree[relpath] = string(body)
		}
		return nil
	}); err != nil {
		return "", nil
	}
	data, err := json.Marshal(tree)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.WithPadding(base64.StdPadding).EncodeToString(data), nil
}

// TreeToDir restore to docker project
func TreeToDir(tree map[string]string, outputDir string) error {
	for k, v := range tree {
		fn := filepath.Join(outputDir, k)
		if err := utils.EnsureFile(fn); err != nil {
			return err
		}
		if err := ioutil.WriteFile(fn, []byte(v), 0600); err != nil {
			return err
		}
	}
	return nil
}
