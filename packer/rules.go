package packer

import (
	"fmt"
	"path/filepath"
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
