package packer

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/metrue/fx/utils"
)

func TestPacker(t *testing.T) {
	t.Run("Pack directory with Dockerfile in it", func(t *testing.T) {
		input := "./fixture/p1"
		output := "output-1"
		defer func() {
			os.RemoveAll(output)
		}()
		if err := Pack(output, input); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Pack directory only fx.js in it", func(t *testing.T) {
		input := "./fixture/p2"
		output := "output-2"
		defer func() {
			os.RemoveAll(output)
		}()
		if err := Pack(output, input); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Pack directory with fx.js and helper in it", func(t *testing.T) {
		input := "./fixture/p3"
		output := "output-3"
		defer func() {
			os.RemoveAll(output)
		}()
		if err := Pack(output, input); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Pack files list with fx.js in it", func(t *testing.T) {
		handleFile := "./fixture/p3/fx.js"
		helperFile := "./fixture/p3/helper.js"
		output := "output-4"
		defer func() {
			os.RemoveAll(output)
		}()
		if err := Pack(output, handleFile, helperFile); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Pack files list without fx.js in it", func(t *testing.T) {
		f1 := "./fixture/p3/helper.js"
		f2 := "./fixture/p3/helper.js"
		output := "output-5"
		defer func() {
			os.RemoveAll(output)
		}()
		if err := Pack(output, f1, f2); err == nil {
			t.Fatalf("should report error when there is not Dockerfile or fx.[ext] in it")
		}
	})
}

func TestTreeAndUnTree(t *testing.T) {
	_, err := PackIntoK8SConfigMapFile("./fixture/p1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGenerate(t *testing.T) {
	langs := []string{
		"d",
		"go",
		"java",
		"julia",
		"node",
		"php",
		"python",
		"ruby",
		"rust",
	}
	for _, lang := range langs {
		output := fmt.Sprintf("output-%s-%d", lang, time.Now().Unix())
		defer func() {
			os.RemoveAll(output)
		}()
		if err := restore(output, lang); err != nil {
			t.Fatal(err)
		}
		diffCmd := exec.Command("diff", "-r", output, "./images/"+lang)
		if stdoutStderr, err := diffCmd.CombinedOutput(); err != nil {
			fmt.Printf("%s\n", stdoutStderr)
			t.Fatal(err)
		}
	}
}

func TestMerge(t *testing.T) {
	// TODO should check the merge result
	t.Run("NoInput", func(t *testing.T) {
		dest := "./dest"
		_ = utils.EnsureDir("./dest")
		defer func() {
			os.RemoveAll(dest)
		}()

		if err := merge(dest); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Files", func(t *testing.T) {
		dest := "./dest"
		_ = utils.EnsureDir("./dest")
		defer func() {
			os.RemoveAll(dest)
		}()

		f1, err := ioutil.TempFile("", "fx.*.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(f1.Name())

		f2, err := ioutil.TempFile("", "fx.*.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(f2.Name())

		if err := merge(dest, f1.Name(), f2.Name()); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Directories", func(t *testing.T) {
		dest := "./dest"
		_ = utils.EnsureDir("./dest")
		defer func() {
			os.RemoveAll(dest)
		}()

		if err := merge(dest, "./fixture/p1"); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Files and Directories", func(t *testing.T) {
		dest := "./dest"
		_ = utils.EnsureDir("./dest")
		defer func() {
			os.RemoveAll(dest)
		}()

		f1, err := ioutil.TempFile("", "fx.*.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(f1.Name())

		f2, err := ioutil.TempFile("", "fx.*.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(f2.Name())

		if err := merge(dest, "./fixture/p1", f1.Name(), f2.Name()); err != nil {
			t.Fatal(err)
		}
	})
}
