package hook

import (
	"os"
	"path/filepath"

	"github.com/metrue/fx/utils"
)

// HookNameBeforeBuild before build hook
const HookNameBeforeBuild = "before_build"

// RunBeforeBuildHook trigger before_build hook
func RunBeforeBuildHook(workdir string) error {
	hooks, err := descovery("")
	if err != nil {
		return err
	}
	for _, h := range hooks {
		if h.Name() == HookNameBeforeBuild {
			if err := h.Run(workdir); err != nil {
				return err
			}
		}
	}
	return nil
}

func descovery(hookdir string) ([]*Hook, error) {
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
