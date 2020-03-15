package julia

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/metrue/fx/utils"
)

func TestJuliaBundler(t *testing.T) {
	t.Run("Scaffold", func(t *testing.T) {
		outputDir, err := ioutil.TempDir("", "fx_julia")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(outputDir)

		julia := New()
		if err := julia.Scaffold(outputDir); err != nil {
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
		fd, err := ioutil.TempFile("", "fx_func_*.julia")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(fd.Name())

		content := `struct Input
    a::Number
    b::Number
end

fx = function(input::Input)
    return input.a + input.b
end
		`
		err = ioutil.WriteFile(fd.Name(), []byte(content), 0666)
		if err != nil {
			t.Fatal(err)
		}

		outputDir, err := ioutil.TempDir("", "fx_julia")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(outputDir)

		julia := New()
		if err := julia.Bundle(outputDir, fd.Name()); err != nil {
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

		preHandleFunc, err := ioutil.ReadFile("./assets/fx.jl")
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
		fd, err := ioutil.TempFile("", "fx_func_*.js")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(fd.Name())

		content, err := ioutil.ReadFile("./assets/fx.jl")
		if err != nil {
			t.Fatal(err)
		}
		err = ioutil.WriteFile(fd.Name(), content, 0666)
		if err != nil {
			t.Fatal(err)
		}

		addFunc := `
module.exports = (a, b) => a+b
  `
		addFd, err := ioutil.TempFile("", "fx_add_func_*.js")
		if err != nil {
			t.Fatal(err)
		}

		err = ioutil.WriteFile(addFd.Name(), []byte(addFunc), 0644)
		if err != nil {
			t.Fatal(err)
		}

		outputDir, err := ioutil.TempDir("", "fx_julia")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(outputDir)

		julia := New()
		if err := julia.Bundle(outputDir, fd.Name(), addFd.Name()); err != nil {
			t.Fatal(err)
		}

		diff, pre, cur, err := utils.Diff("./assets", outputDir)
		if err != nil {
			t.Fatal(err)
		}

		if !diff {
			t.Fatalf("handle functino should be changed")
		}

		if !reflect.DeepEqual(cur, []byte(addFunc)) {
			t.Fatalf("it should be %s but got %s", content, cur)
		}
		if pre != nil {
			t.Fatal(pre)
		}
	})
}
