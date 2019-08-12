package packer

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/metrue/fx/types"
)

// DockerPacker pack a function source code to a Docker build-able project
type DockerPacker struct {
	box packr.Box
}

func isHandler(lang string, name string) bool {
	basename := filepath.Base(name)
	nameWithoutExt := strings.TrimSuffix(basename, filepath.Ext(basename))
	return nameWithoutExt == "fx" ||
		nameWithoutExt == "Fx" || // Fx is for Java
		nameWithoutExt == "mod" // mod.rs is for Rust
}

// NewDockerPacker new a Docker packer
func NewDockerPacker(box packr.Box) *DockerPacker {
	return &DockerPacker{box: box}
}

// Pack pack a single function source code to be project
func (p *DockerPacker) Pack(serviceName string, fn types.ServiceFunctionSource) (types.Project, error) {
	var files []types.ProjectSourceFile
	for _, name := range p.box.List() {
		prefix := fmt.Sprintf("%s/", fn.Language)
		if strings.HasPrefix(name, prefix) {
			content, err := p.box.FindString(name)
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
