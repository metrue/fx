package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/otiai10/copy"
)

// Merge inputs into dest
func Merge(dest string, input ...string) error {
	for _, file := range input {
		stat, err := os.Stat(file)
		if err != nil {
			return err
		}
		if stat.Mode().IsRegular() {
			targetFilePath := filepath.Join(dest, stat.Name())
			if err := EnsureFile(targetFilePath); err != nil {
				return err
			}
			body, err := ioutil.ReadFile(file)
			if err != nil {
				return err
			}
			if err := ioutil.WriteFile(targetFilePath, body, 0600); err != nil {
				return err
			}
		} else if stat.Mode().IsDir() {
			if err := filepath.Walk(file, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}

				stat, err := os.Stat(path)
				if err != nil {
					return err
				}
				if stat.Mode().IsRegular() {
					destDir := filepath.Join(dest, filepath.Dir(path))
					if err := EnsureDir(destDir); err != nil {
						return err
					}

					if err := copy.Copy(file, destDir); err != nil {
						return err
					}
				}

				return nil
			}); err != nil {
				return err
			}
		}
	}
	return nil
}

// Diff two directory
func Diff(src string, dst string) (diff bool, pre []byte, cur []byte, err error) {
	src, err = filepath.Abs(src)
	if err != nil {
		return true, nil, nil, err
	}

	dst, err = filepath.Abs(dst)
	if err != nil {
		return true, nil, nil, err
	}

	srcMap := map[string][]byte{}
	if err := filepath.Walk(src, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		if IsDir(path) {
			return nil
		}
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		p := strings.Replace(path, src, "", -1)
		srcMap[p] = content
		return nil
	}); err != nil {
		return true, nil, nil, err
	}

	dstMap := map[string][]byte{}
	if err := filepath.Walk(dst, func(path string, info os.FileInfo, e error) error {
		if e != nil {
			return e
		}
		if IsDir(path) {
			return nil
		}
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		p := strings.Replace(path, dst, "", -1)
		dstMap[p] = content
		return nil
	}); err != nil {
		return true, nil, nil, err
	}

	for k, v := range srcMap {
		if !reflect.DeepEqual(dstMap[k], v) {
			return true, v, dstMap[k], nil
		}
	}

	for k, v := range dstMap {
		if !reflect.DeepEqual(srcMap[k], v) {
			return true, srcMap[k], v, nil
		}
	}

	return false, nil, nil, nil
}
