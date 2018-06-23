package utils

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

// Download a resource from URL to given path
func Download(filepath string, url string) (err error) {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

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
		path := filepath.Join(target, file.Name)
		if file.FileInfo().IsDir() {
			os.MkdirAll(path, file.Mode())
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

// Create Dir if not exist
func EnsurerDir(dir string) (err error) {
	if _, statError := os.Stat(dir); os.IsNotExist(statError) {
		mkError := os.MkdirAll(dir, os.ModePerm)
		return mkError
	}
	return nil
}

// TarDir builds a tar from directory
func TarDir(srcDir string, desFileName string) {
	dir, err := os.Open(srcDir)
	if err != nil {
		panic(err)
	}
	defer dir.Close()

	files, err := dir.Readdir(0)
	if err != nil {
		panic(err)
	}

	// create tar file
	tarfile, err := os.Create(desFileName)
	if err != nil {
		panic(err)
	}
	defer tarfile.Close()

	var fileWriter io.WriteCloser = tarfile

	tarfileWriter := tar.NewWriter(fileWriter)
	defer tarfileWriter.Close()

	for _, fileInfo := range files {
		if fileInfo.IsDir() {
			continue
		}

		file, err := os.Open(dir.Name() + string(filepath.Separator) + fileInfo.Name())
		checkerror(err)
		defer file.Close()

		// prepare the tar header
		header := new(tar.Header)
		header.Name = file.Name()
		header.Size = fileInfo.Size()
		header.Mode = int64(fileInfo.Mode())
		header.ModTime = fileInfo.ModTime()

		err = tarfileWriter.WriteHeader(header)
		checkerror(err)

		_, err = io.Copy(tarfileWriter, file)
		checkerror(err)
	}
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

func checkerror(err error) {
	if err != nil {
		panic(err)
	}
}

// GetCurrentExecPath parses a path from running executable/go file
func GetCurrentExecPath() (scriptPath string) {
	scriptPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	return scriptPath
}

func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}

// GetHostIP returns the host's IP
func GetHostIP() (ip net.IP) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	HandleError(err)

	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP
}

// get programming language from file name extension
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
	}
	return extLangMap[filepath.Ext(fileName)]
}
