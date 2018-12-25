package bundler

import (
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestBundle(t *testing.T) {
	targetDir, _ := ioutil.TempDir("", "image_test")
	Bundle(targetDir, "go", []byte("import \"fmt\""))

	files, _ := ioutil.ReadDir(targetDir)
	if len(files) != 3 {
		t.Errorf("files number not correct, got: %d, want: %d.", len(files), 3)
	}

	if files[0].Name() != "Dockerfile" {
		t.Errorf("Dockerfile not correct, got: %s, want: %s.", files[0], "Dockerfile")
	}

	if files[1].Name() != "app.go" {
		t.Errorf("app.go not correct, got: %s, want: %s.", files[1], "app.go")
	}

	if files[2].Name() != "fx.go" {
		t.Errorf("fx.go not correct, got: %s, want: %s.", files[2], "fx.go")
	}

	filePath := path.Join(targetDir, files[2].Name())
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Errorf("open fx.go error: %s", err)
	}

	if string(data) != "import \"fmt\"" {
		t.Errorf("content of fx.go not correct, got: %s, want: %s.", data, "import \"fmt\"")
	}

	defer os.RemoveAll(targetDir)
}
