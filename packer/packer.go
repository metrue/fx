package packer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/metrue/fx/utils"
	"github.com/otiai10/copy"
)

var presets packr.Box

func init() {
	presets = packr.NewBox("./images")
}

// Pack pack a file or directory into a Docker project
func Pack(output string, input ...string) error {
	if len(input) == 0 {
		return fmt.Errorf("source file or directory required")
	}

	var lang string
	for _, f := range input {
		if utils.IsRegularFile(f) {
			lang = langFromFileName(f)
		} else if utils.IsDir(f) {
			if err := filepath.Walk(f, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				if utils.IsRegularFile(path) {
					lang = langFromFileName(path)
				}
				return nil
			}); err != nil {
				return err
			}
		}
	}

	if lang == "" {
		return fmt.Errorf("could not tell programe language of your input source codes")
	}

	if err := restore(output, lang); err != nil {
		return err
	}

	if len(input) == 1 {
		stat, err := os.Stat(input[0])
		if err != nil {
			return err
		}
		if stat.Mode().IsRegular() {
			if err := filepath.Walk(output, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if isHandler(path) {
					if err := copy.Copy(input[0], path); err != nil {
						return err
					}
				}

				return nil
			}); err != nil {
				return err
			}
		}
		return nil
	}

	if !hasFxHandleFile(input...) {
		msg := `it requires a fx handle file when input is not a single file function, e.g.  
fx.go for Golang
Fx.java for Java
fx.php for PHP
fx.py for Python
fx.js for JavaScript or Node
fx.rb for Ruby
fx.jl for Julia
fx.d for D`
		return fmt.Errorf(msg)
	}

	if err := merge(output, input...); err != nil {
		return err
	}
	return nil
}

func restore(output string, lang string) error {
	for _, name := range presets.List() {
		prefix := fmt.Sprintf("%s/", lang)
		if strings.HasPrefix(name, prefix) {
			content, err := presets.FindString(name)
			if err != nil {
				return err
			}

			filePath := filepath.Join(output, strings.Replace(name, prefix, "", 1))
			if err := utils.EnsureFile(filePath); err != nil {
				return err
			}
			if err := ioutil.WriteFile(filePath, []byte(content), 0666); err != nil {
				return err
			}
		}
	}
	return nil
}

func merge(dest string, input ...string) error {
	for _, file := range input {
		stat, err := os.Stat(file)
		if err != nil {
			return err
		}
		if stat.Mode().IsRegular() {
			targetFilePath := filepath.Join(dest, stat.Name())
			if err := utils.EnsureFile(targetFilePath); err != nil {
				return err
			}
			body, err := ioutil.ReadFile(file)
			if err != nil {
				return err
			}
			if err := ioutil.WriteFile(targetFilePath, body, 0644); err != nil {
				return err
			}
		} else if stat.Mode().IsDir() {
			if err := filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if err := copy.Copy(file, dest); err != nil {
					return err
				}
				return nil
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

func langFromFileName(fileName string) string {
	extLangMap := map[string]string{
		".js":   "node",
		".go":   "go",
		".rb":   "ruby",
		".py":   "python",
		".php":  "php",
		".jl":   "julia",
		".java": "java",
		".d":    "d",
		".rs":   "rust",
	}
	return extLangMap[filepath.Ext(fileName)]
}

func hasFxHandleFile(input ...string) bool {
	var handleFile string
	for _, file := range input {
		if utils.IsRegularFile(file) && isHandler(file) {
			handleFile = file
			break
		} else if utils.IsDir(file) {
			if err := filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if utils.IsRegularFile(path) && isHandler(info.Name()) {
					handleFile = path
				}
				return nil
			}); err != nil {
				return false
			}
		}
	}

	return handleFile != ""
}

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
		if err := ioutil.WriteFile(fn, []byte(v), 0666); err != nil {
			return err
		}
	}
	return nil
}
