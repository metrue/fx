package packer

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestDeps(t *testing.T) {
	cases := []struct {
		name    string
		content string
		expect  []string
	}{
		{name: "import with singe quote", content: "import 'pkg1'", expect: []string{"pkg1"}},
		{name: "import with double quote", content: `import "pkg1"`, expect: []string{"pkg1"}},
		{name: "require with singe quote", content: "require 'pkg1'", expect: []string{"pkg1"}},
		{name: "require with double quote", content: `require "pkg1"`, expect: []string{"pkg1"}},
		{
			name: "require and import with singe quote",
			content: `
require 'pkg1'
import 'pkg2'`,
			expect: []string{"pkg1", "pkg2"},
		},
		{
			name: "require and import with double quote",
			content: `
require "pkg1"
import "pkg2"`,
			expect: []string{"pkg1", "pkg2"},
		},
		{
			name: "require and import with double quote and single quote",
			content: `
require 'pkg1'
import "pkg2"`,
			expect: []string{"pkg1", "pkg2"},
		},
	}

	for _, c := range cases {
		fd, err := ioutil.TempFile("", "dep_pkg_*.js")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(fd.Name())

		err = ioutil.WriteFile(fd.Name(), []byte(c.content), 0666)
		if err != nil {
			t.Fatal(err)
		}

		deps, err := Deps(fd.Name())
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(deps, c.expect) {
			t.Fatalf("%s: should get %s but got %s", c.name, c.expect, deps)
		}
	}
}
