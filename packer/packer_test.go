package packer

import (
	"os"
	"testing"

	"github.com/metrue/fx/types"
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

	t.Run("pack", func(t *testing.T) {
		mockSource := `
module.exports = ({a, b}) => {
	return a + b
}
`
		fn := types.Func{
			Language: "node",
			Source:   mockSource,
		}

		serviceName := "service-mock"
		project, err := pack(serviceName, fn)
		if err != nil {
			t.Fatal(err)
		}

		if project.Name != serviceName {
			t.Fatalf("should get %s but got %s", serviceName, project.Name)
		}

		if project.Language != "node" {
			t.Fatal("incorrect Language")
		}

		if len(project.Files) != 3 {
			t.Fatal("node project should have 3 files")
		}

		for _, file := range project.Files {
			if file.Path == "fx.js" {
				if file.IsHandler == false {
					t.Fatal("fx.js should be handler")
				}
				if file.Body != mockSource {
					t.Fatalf("should get %s but got %v", mockSource, file.Body)
				}
			} else if file.Path == "Dockerfile" {
				if file.IsHandler == true {
					t.Fatalf("should get %v but got %v", false, file.IsHandler)
				}
			} else {
				if file.IsHandler == true {
					t.Fatalf("should get %v but %v", false, file.IsHandler)
				}
			}
		}
	})
}

func TestTreeAndUnTree(t *testing.T) {
	_, err := PackIntoK8SConfigMapFile("./fixture/p1")
	if err != nil {
		t.Fatal(err)
	}
}
