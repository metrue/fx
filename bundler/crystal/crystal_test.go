package crystal

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/metrue/fx/utils"
)

func TestCrystalBundler(t *testing.T) {
	t.Run("Scaffold", func(t *testing.T) {
		outputDir, err := ioutil.TempDir("", "fx_crystal")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(outputDir)

		crystal := New()
		if err := crystal.Scaffold(outputDir); err != nil {
			t.Fatal(err)
		}

		diff, _, _, err := utils.Diff(outputDir, "./assets")
		if err != nil {
			t.Fatal(err)
		}

		if diff {
			t.Fatalf("%s is not equal with %s", outputDir, "./assets")
		}
	})

	t.Run("BundleSingleFunc", func(t *testing.T) {
		fd, err := ioutil.TempFile("", "fx_func_*.crystal")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(fd.Name())

		content := `def fx(ctx)
	"hello world"
end`
		err = ioutil.WriteFile(fd.Name(), []byte(content), 0666)
		if err != nil {
			t.Fatal(err)
		}

		outputDir, err := ioutil.TempDir("", "fx_crystal")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(outputDir)

		crystal := New()
		if err := crystal.Bundle(outputDir, fd.Name()); err != nil {
			t.Fatal(err)
		}

		diff, pre, cur, err := utils.Diff("./assets", outputDir)
		if err != nil {
			t.Fatal(err)
		}

		if !diff {
			t.Fatalf("handle function should be changed")
		}

		if !reflect.DeepEqual(cur, []byte(content)) {
			t.Fatalf("it should be %s but got %s", content, cur)
		}

		preHandleFunc, err := ioutil.ReadFile("./assets/fx.cr")
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(pre, preHandleFunc) {
			{
			}
			t.Fatalf("it should get %s but got %s", preHandleFunc, pre)
		}
	})

	t.Run("BundleFuncAndDeps", func(t *testing.T) {
		fd, err := ioutil.TempFile("", "fx_func_*.cr")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(fd.Name())

		content, err := ioutil.ReadFile("./assets/fx.cr")
		if err != nil {
			t.Fatal(err)
		}
		err = ioutil.WriteFile(fd.Name(), content, 0666)
		if err != nil {
			t.Fatal(err)
		}

		addFunc := `def fx(a, b)
	a + b
end
`
		addFd, err := ioutil.TempFile("", "fx_add_func_*.cr")
		if err != nil {
			t.Fatal(err)
		}

		err = ioutil.WriteFile(addFd.Name(), []byte(addFunc), 0644)
		if err != nil {
			t.Fatal(err)
		}

		outputDir, err := ioutil.TempDir("", "fx_crystal")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(outputDir)

		crystal := New()
		if err := crystal.Bundle(outputDir, fd.Name(), addFd.Name()); err != nil {
			t.Fatal(err)
		}

		diff, pre, cur, err := utils.Diff("./assets", outputDir)
		if err != nil {
			t.Fatal(err)
		}

		if !diff {
			t.Fatalf("handle function should be changed")
		}

		if !reflect.DeepEqual(cur, []byte(addFunc)) {
			t.Fatalf("it should be %s but got %s", content, cur)
		}
		if pre != nil {
			t.Fatal(pre)
		}
	})
}
