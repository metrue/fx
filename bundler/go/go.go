package golang

import (
	"github.com/gobuffalo/packr/v2"
	"github.com/metrue/fx/bundler"
)

// Gin defines javascript bundler
type Gin struct {
	assets *packr.Box
}

// New a koa bundler
func New() *Gin {
	return &Gin{
		assets: packr.New("", "./assets"),
	}
}

// Scaffold a koa app
func (g *Gin) Scaffold(output string) error {
	return bundler.Restore(g.assets, output)
}

// Bundle a function into a koa project
func (g *Gin) Bundle(output string, fn string, deps ...string) error {
	return bundler.Bundle(g.assets, output, "go", fn, deps...)
}

var (
	_ bundler.Bundler = &Gin{}
)
