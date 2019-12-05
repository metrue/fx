package packer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"encoding/base64"
	"encoding/json"

	"github.com/gobuffalo/packr"
	"github.com/metrue/fx/types"
	"github.com/metrue/fx/utils"
	"github.com/otiai10/copy"
	"github.com/pkg/errors"
)

// Pack pack a file or directory into a Docker project
func Pack(output string, input ...string) error {
	if len(input) == 0 {
		return fmt.Errorf("source file or directory required")
	}

	if len(input) == 1 {
		file := input[0]
		stat, err := os.Stat(file)
		if err != nil {
			return err
		}
		if !stat.IsDir() {
			lang := utils.GetLangFromFileName(file)
			body, err := ioutil.ReadFile(file)
			if err != nil {
				return errors.Wrap(err, "read source failed")
			}
			fn := types.Func{
				Language: lang,
				Source:   string(body),
			}
			if err := PackIntoDir(fn, output); err != nil {
				return err
			}
			return nil
		}
	}

	workdir := fmt.Sprintf("./fx-%d", time.Now().Unix())
	defer os.RemoveAll(workdir)

	for _, f := range input {
		if err := copy.Copy(f, filepath.Join(workdir, f)); err != nil {
			return err
		}
	}

	if dockerfile, has := hasDockerfileInDir(workdir); has {
		return copy.Copy(filepath.Dir(dockerfile), output)
	}

	if f, has := hasFxHandleFileInDir(workdir); has {
		lang := utils.GetLangFromFileName(f)
		body, err := ioutil.ReadFile(f)
		if err != nil {
			return errors.Wrap(err, "read source failed")
		}
		fn := types.Func{
			Language: lang,
			Source:   string(body),
		}
		if err := PackIntoDir(fn, output); err != nil {
			return err
		}
		return copy.Copy(filepath.Dir(f), output)
	}

	return fmt.Errorf("input directories or files has no Dockerfile or file with fx as name, e.g. fx.js")
}

func hasDockerfileInDir(dir string) (string, bool) {
	var dockerfile string
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// nolint
		if !info.IsDir() && info.Name() == "Dockerfile" {
			dockerfile = path
		}
		return nil
	}); err != nil {
		return "", false
	}
	if dockerfile == "" {
		return "", false
	}
	return dockerfile, true
}

func hasFxHandleFileInDir(dir string) (string, bool) {
	var handleFile string
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && isHandler(info.Name()) {
			handleFile = path
		}
		return nil
	}); err != nil {
		return "", false
	}
	if handleFile == "" {
		return "", false
	}
	return handleFile, true
}

// Pack a function to be a docker project which is web service, handle the imcome request with given function
func pack(svcName string, fn types.Func) (types.Project, error) {
	box := packr.NewBox("./images")
	pkr := NewDockerPacker(box)
	return pkr.Pack(svcName, fn)
}

// PackIntoDir pack service code into directory
func PackIntoDir(fn types.Func, outputDir string) error {
	project, err := pack("", fn)
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
func PackIntoK8SConfigMapFile(dir string) (string, error) {
	tree := map[string]string{}
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relpath := strings.Replace(path, dir, "", 1)
			body, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			tree[relpath] = string(body)
		}
		return nil
	}); err != nil {
		return "", nil
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
