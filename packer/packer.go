package packer

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"encoding/base64"
	"encoding/json"

	"github.com/gobuffalo/packr"
	"github.com/metrue/fx/types"
	"github.com/metrue/fx/utils"
)

// Packer interface
type Packer interface {
	Pack(serviceName string, fn types.Func) (types.Project, error)
}

// Pack a function to be a docker project which is web service, handle the imcome request with given function
func Pack(svcName string, fn types.Func) (types.Project, error) {
	box := packr.NewBox("./images")
	pkr := NewDockerPacker(box)
	return pkr.Pack(svcName, fn)
}

// PackIntoDir pack service code into directory
func PackIntoDir(fn types.Func, outputDir string) error {
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

// PackIntoK8SConfigMapFile pack function a K8S config map file
func PackIntoK8SConfigMapFile(fn types.Func) (string, error) {
	project, err := Pack("", fn)
	if err != nil {
		return "", err
	}
	tree := map[string]string{}
	for _, file := range project.Files {
		tree[file.Path] = file.Body
	}

	data, err := json.Marshal(tree)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.WithPadding(base64.StdPadding).EncodeToString(data), nil
}

// TreeToDir restore to docker project
func TreeToDir(tree map[string]string, outputDir string) error {
	for k, v := range tree {
		fn := filepath.Join(outputDir, k)
		if err := utils.EnsureFile(fn); err != nil {
			return err
		}
		if err := ioutil.WriteFile(fn, []byte(v), 0666); err != nil {
			return err
		}
	}
	return nil
}

// PackIntoTar pack service code into directory
func PackIntoTar(fn types.Func, path string) error {
	tarDir, err := ioutil.TempDir("/tmp", "fx-tar")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tarDir)

	if err := PackIntoDir(fn, tarDir); err != nil {
		return err
	}

	return utils.TarDir(tarDir, path)
}
