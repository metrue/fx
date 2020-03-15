package java

import (
	"github.com/gobuffalo/packr/v2"
	"github.com/metrue/fx/bundler"
)

// Java defines javascript bundler
type Java struct {
	assets *packr.Box
}

// New a koa bundler
func New() *Java {
	return &Java{
		assets: packr.New("java", "./assets"),
	}
}

// Scaffold a koa app
func (k *Java) Scaffold(output string) error {
	return bundler.Restore(k.assets, output)
}

// Bundle a function into a koa project
func (k *Java) Bundle(output string, fn string, deps ...string) error {
	return bundler.Bundle(k.assets, output, "java", fn, deps...)
}

var (
	_ bundler.Bundler = &Java{}
)
