package hook

import (
	"os"
	"path/filepath"

	"github.com/metrue/fx/utils"
)

// HookNameBeforeBuild before build hook
const HookNameBeforeBuild = "before_build"

// Descovery hooks in given
func Descovery(hookdir string) ([]*Hook, error) {
	if hookdir == "" {
		dir, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		hookdir = filepath.Join(dir, ".hooks")
	}

	hooks := []*Hook{}
	if !utils.IsDir(hookdir) {
		return hooks, nil
	}
	if err := filepath.Walk(hookdir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == HookNameBeforeBuild {
			hooks = append(hooks, New("before_build", path, ""))
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return hooks, nil
}
