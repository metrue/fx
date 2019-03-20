package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mholt/archiver"
)

// TarDir make tar of a dir
func TarDir(dir string, tarFileName string) error {
	if !strings.HasSuffix(tarFileName, "tar") {
		return fmt.Errorf("destination file name must end with .tar")
	}

	if err := os.Chdir(dir); err != nil {
		return err
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	files, err := ioutil.ReadDir(wd)
	if err != nil {
		return err
	}

	sources := []string{}
	for _, file := range files {
		sources = append(sources, file.Name())
	}

	return archiver.Archive(sources, tarFileName)
}
