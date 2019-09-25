package packer

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/gobuffalo/packr"
	"github.com/metrue/fx/types"
	"github.com/metrue/fx/utils"
)

// Packer interface
type Packer interface {
	Pack(serviceName string, fn types.ServiceFunctionSource) (types.Project, error)
}

// Pack a function to be a docker project which is web service, handle the imcome request with given function
func Pack(svcName string, fn types.ServiceFunctionSource) (types.Project, error) {
	box := packr.NewBox("./images")
	pkr := NewDockerPacker(box)
	return pkr.Pack(svcName, fn)
}

// PackIntoDir pack service code into directory
func PackIntoDir(lang string, source string, outputDir string) error {
	fn := types.ServiceFunctionSource{
		Language: lang,
		Source:   source,
	}
	project, err := Pack("", fn)
	if err != nil {
		return err
	}
	for _, file := range project.Files {
		tmpfn := filepath.Join(outputDir, file.Path)
		if err := utils.EnsureFile(tmpfn); err != nil {
			return err
		}
		if err := ioutil.WriteFile(tmpfn, []byte(file.Body), 0666); err != nil {
			return err
		}
	}
	return nil
}

// PackIntoTar pack service code into directory
func PackIntoTar(lang string, source string, path string) error {
	tarDir, err := ioutil.TempDir("/tmp", "fx-tar")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tarDir)

	if err := PackIntoDir(lang, source, tarDir); err != nil {
		return err
	}

	return utils.TarDir(tarDir, path)
}
