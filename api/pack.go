package api

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/metrue/fx/types"
)

func isHandler(lang string, name string) bool {
	basename := filepath.Base(name)
	nameWithoutExt := strings.TrimSuffix(basename, filepath.Ext(basename))
	return nameWithoutExt == "fx" ||
		nameWithoutExt == "Fx" || // Fx is for Java
		nameWithoutExt == "mod" // mod.rs is for Rust
}

// Pack pack a single function source code to be project
func (api *API) Pack(serviceName string, fn types.ServiceFunctionSource) (types.Project, error) {
	var files []types.ProjectSourceFile
	for _, name := range api.box.List() {
		prefix := fmt.Sprintf("%s/", fn.Language)
		if strings.HasPrefix(name, prefix) {
			content, err := api.box.FindString(name)
			if err != nil {
				return types.Project{}, err
			}

			// if preset's file is handler function of project, replace it with give one
			if isHandler(fn.Language, name) {
				files = append(files, types.ProjectSourceFile{
					Path:      strings.Replace(name, prefix, "", 1),
					Body:      fn.Source,
					IsHandler: true,
				})
			} else {
				files = append(files, types.ProjectSourceFile{
					Path:      strings.Replace(name, prefix, "", 1),
					Body:      content,
					IsHandler: false,
				})
			}
		}
	}

	return types.Project{
		Name:     serviceName,
		Files:    files,
		Language: fn.Language,
	}, nil
}
