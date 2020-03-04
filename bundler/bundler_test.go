package bundler

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

func merge(dumy ...string) error {
	return nil
}

func TestBundler(t *testing.T) {
	t.Run("Restore", func(t *testing.T) {
		t.Skip()
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
			if err := Restore(nil, output); err != nil {
				t.Fatal(err)
			}
			diffCmd := exec.Command("diff", "-r", output, "./images/"+lang)
			if stdoutStderr, err := diffCmd.CombinedOutput(); err != nil {
				fmt.Printf("%s\n", stdoutStderr)
				t.Fatal(err)
			}
		}
	})

	t.Run("Bundle", func(t *testing.T) {
		t.Skip()
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
	})
}
