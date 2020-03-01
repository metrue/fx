package packer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// JavaScriptExt extension of JavaScript source file
const JavaScriptExt = ".js"

func depsJavaScript(src string) ([]string, error) {
	deps := []string{}
	file, err := os.Open(src)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "import") || strings.HasPrefix(line, "require") {
			reg, err := regexp.Compile(`(\'.*\'|\".*\")`)
			if err != nil {
				return nil, err
			}
			str := reg.FindString(line)
			if str != "" {
				deps = append(deps, strings.ReplaceAll(strings.ReplaceAll(str, "\"", ""), "'", ""))
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return deps, nil
}

// Deps to get the dependencies required from a source code in a file
func Deps(src string) ([]string, error) {
	switch filepath.Ext(src) {
	case JavaScriptExt:
		{
			return depsJavaScript(src)
		}
	}
	return nil, fmt.Errorf("not support language: %s", filepath.Ext(src))
}
