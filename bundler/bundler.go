package bundler

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobuffalo/packd"
	"github.com/gobuffalo/packr/v2"
	"github.com/metrue/fx/constants"
	"github.com/metrue/fx/utils"
)

// Bundler defines interface
type Bundler interface {
	Scaffold(output string) error
	Bundle(output string, fn string, deps ...string) error
}

// IsHandler check if it's handle file
func IsHandler(name string, lang string) bool {
	basename := filepath.Base(name)
	nameWithoutExt := strings.TrimSuffix(basename, filepath.Ext(basename))
	if constants.ExtLangMapping[filepath.Ext(basename)] != lang {
		return false
	}

	return (nameWithoutExt == "fx" ||
		// Fx is for Java
		nameWithoutExt == "Fx" ||
		// mod.rs is for Rust)
		nameWithoutExt == "mod")
}

// Restore directory from packr box
func Restore(box *packr.Box, output string) error {
	if err := box.Walk(func(name string, fd packd.File) error {
		content, err := box.Find(name)
		if err != nil {
			return err
		}

		dest := filepath.Join(output, name)
		if err := utils.EnsureFile(dest); err != nil {
			return err
		}

		if err := ioutil.WriteFile(dest, content, 0644); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}

// Bundle bundle a function
func Bundle(box *packr.Box, output string, language string, fn string, deps ...string) error {
	if err := Restore(box, output); err != nil {
		return err
	}

	if err := utils.Merge(output, deps...); err != nil {
		return err
	}

	// Replace the default handler source file with given function source file
	if err := filepath.Walk(output, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if IsHandler(path, language) {
			if err := utils.CopyFile(fn, path); err != nil {
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
