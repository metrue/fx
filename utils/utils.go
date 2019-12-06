package utils

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Download a resource from URL to given path
func Download(filepath string, url string) (err error) {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// nolint: gosec
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

// Unzip a folder to destination
func Unzip(source string, target string) (err error) {
	reader, err := zip.OpenReader(source)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(target, 0755); err != nil {
		return err
	}

	for _, file := range reader.File {
		//nolint: gosec
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			if err := os.MkdirAll(path, file.Mode()); err != nil {
				return err
			}
			continue
		}

		fileReader, err := file.Open()
		if err != nil {
			return err
		}
		defer fileReader.Close()

		targetFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer targetFile.Close()

		if _, err := io.Copy(targetFile, fileReader); err != nil {
			return err
		}
	}

	return nil
}

// CopyFile from src to dst
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

// EnsureDir Create Dir if not exist
func EnsureDir(dir string) (err error) {
	if _, statError := os.Stat(dir); os.IsNotExist(statError) {
		mkError := os.MkdirAll(dir, os.ModePerm)
		return mkError
	}
	return nil
}

// EnsureFile ensure a file
func EnsureFile(fullpath string) error {
	dir := path.Dir(fullpath)
	err := EnsureDir(dir)
	if err != nil {
		return err
	}
	_, err = os.OpenFile(fullpath, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	return nil
}

// IsDir if given path is a directory
func IsDir(dir string) bool {
	stat, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

// IsRegularFile if given path is a regular
func IsRegularFile(file string) bool {
	stat, err := os.Stat(file)
	if err != nil {
		return false
	}
	return stat.Mode().IsRegular()
}

// IsPathExists checks whether a path exists or if failed to check
func IsPathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// GetCurrentExecPath parses a path from running executable/go file
func GetCurrentExecPath() (scriptPath string) {
	scriptPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return scriptPath
}

// GetHostIP returns the host's IP
func GetHostIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return net.IP{}, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}

// GetLangFromFileName get programming language from file name extension
func GetLangFromFileName(fileName string) (lang string) {
	extLangMap := map[string]string{
		".js":   "node",
		".go":   "go",
		".rb":   "ruby",
		".py":   "python",
		".php":  "php",
		".jl":   "julia",
		".java": "java",
		".d":    "d",
		".rs":   "rust",
	}
	return extLangMap[filepath.Ext(fileName)]
}

// PairsToParams make "a=1, b=2" to be {"a": "1", "b": "2"}
func PairsToParams(pairs []string) map[string]string {
	params := map[string]string{}
	for _, pair := range pairs {
		subs := strings.Split(pair, "=")
		if len(subs) == 2 {
			params[subs[0]] = subs[1]
		}
	}
	return params
}

// OutputJSON output json
func OutputJSON(v interface{}) error {
	bytes, err := json.MarshalIndent(v, "", "\t")
	if err != nil {
		return fmt.Errorf("could marshal %v : %v", v, err)
	}
	fmt.Println(string(bytes))

	return nil
}
