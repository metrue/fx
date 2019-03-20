package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestTarDir(t *testing.T) {
	tmp, err := ioutil.TempFile("/tmp", "test-tar-dir")
	if err != nil {
		t.Fatal(err)
	}
	os.Remove(tmp.Name())

	err = TarDir(".", fmt.Sprintf("%s.tar", tmp.Name()))
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tmp.Name())
}
