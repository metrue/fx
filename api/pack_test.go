package api

import (
	"testing"

	"github.com/metrue/fx/types"
)

func TestPacker(t *testing.T) {
	api := NewWithDockerRemoteAPI("127.0.0.1:1234", "0.2.1")

	mockSource := `
module.exports = ({a, b}) => {
	return a + b
}
`
	fn := types.ServiceFunctionSource{
		Language: "node",
		Source:   mockSource,
	}

	serviceName := "service-mock"
	project, err := api.Pack(serviceName, fn)
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
