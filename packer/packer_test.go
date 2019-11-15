package packer

import (
	"testing"

	"github.com/metrue/fx/types"
)

func TestPack(t *testing.T) {
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
	project, err := Pack(serviceName, fn)
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
}

func TestTreeAndUnTree(t *testing.T) {
	mockSource := `
package fx;

import org.json.JSONObject;

public class Fx {
    public int handle(JSONObject input) {
        String a = input.get("a").toString();
        String b = input.get("b").toString();
        return Integer.parseInt(a) + Integer.parseInt(b);
    }
}
`
	fn := types.Func{
		Language: "java",
		Source:   mockSource,
	}
	_, err := PackIntoK8SConfigMapFile(fn)
	if err != nil {
		t.Fatal(err)
	}
}
