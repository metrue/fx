package packer

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/metrue/fx/utils"
)

// ExtLangMapping file extension mapping with programming language
var ExtLangMapping = map[string]string{
	".js":   "node",
	".go":   "go",
	".rb":   "ruby",
	".py":   "python",
	".php":  "php",
	".jl":   "julia",
	".java": "java",
	".d":    "d",
	".rs":   "rust",
	".pl":   "perl",
}

func isHandler(name string, lang string) bool {
	basename := filepath.Base(name)
	nameWithoutExt := strings.TrimSuffix(basename, filepath.Ext(basename))
	if ExtLangMapping[filepath.Ext(basename)] != lang {
		return false
	}

	return (nameWithoutExt == "fx" ||
		// Fx is for Java
		nameWithoutExt == "Fx" ||
		// mod.rs is for Rust)
		nameWithoutExt == "mod")
}

func langFromFileName(fileName string) (string, error) {
	if fileName == "" {
		return "", fmt.Errorf("file name should not be empty")
	}

	ext := filepath.Ext(fileName)
	lang, ok := ExtLangMapping[ext]
	if !ok {
		return "", fmt.Errorf("could not find corresponse programming language for file extension %s", ext)
	}
	return lang, nil
}

func hasFxHandleFile(lang string, input ...string) bool {
	var handleFile string
	for _, file := range input {
		if utils.IsRegularFile(file) && isHandler(file, lang) {
			handleFile = file
			break
		} else if utils.IsDir(file) {
			if err := filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				if utils.IsRegularFile(path) && isHandler(info.Name(), lang) {
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
