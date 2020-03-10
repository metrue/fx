package bundle

import (
	"fmt"

	"github.com/metrue/fx/bundler"
	"github.com/metrue/fx/bundler/d"
	golang "github.com/metrue/fx/bundler/go"
	"github.com/metrue/fx/bundler/java"
	"github.com/metrue/fx/bundler/julia"
	"github.com/metrue/fx/bundler/node"
	"github.com/metrue/fx/bundler/perl"
	"github.com/metrue/fx/bundler/python"
	"github.com/metrue/fx/bundler/ruby"
	"github.com/metrue/fx/bundler/rust"
)

// Bundle function to project
func Bundle(workdir string, language string, fn string, deps ...string) error {
	var bundler bundler.Bundler
	switch language {
	case "d":
		bundler = d.New()
	case "node":
		bundler = node.New()
	case "go":
		bundler = golang.New()
	case "java":
		bundler = java.New()
	case "julia":
		bundler = julia.New()
	case "perl":
		bundler = perl.New()
	case "python":
		bundler = python.New()
	case "ruby":
		bundler = ruby.New()
	case "rust":
		bundler = rust.New()
	default:
		return fmt.Errorf("%s not suppported yet", language)
	}
	return bundler.Bundle(workdir, fn, deps...)
}
