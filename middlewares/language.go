package middlewares

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/metrue/fx/context"
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

// Language to find out what language of function is
func Language() func(ctx context.Contexter) (err error) {
	return func(ctx context.Contexter) error {
		sources := ctx.Get("sources").([]string)
		var language string
		for _, f := range sources {
			if utils.IsRegularFile(f) {
				lang, err := langFromFileName(f)
				if err == nil {
					language = lang
				}
			} else if utils.IsDir(f) {
				if err := filepath.Walk(f, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if utils.IsRegularFile(path) {
						lang, err := langFromFileName(path)
						if err == nil {
							language = lang
						}
					}
					return nil
				}); err != nil {
					return err
				}
			}
		}
		if language == "" {
			return fmt.Errorf("could not tell programing language of your source codes")
		}
		ctx.Set("language", language)
		return nil
	}
}
