package d

import (
	"github.com/gobuffalo/packr/v2"
	"github.com/metrue/fx/bundler"
)

// D defines d bundler
type D struct {
	assets *packr.Box
}

// New a koa bundler
func New() *D {
	return &D{
		assets: packr.New("d", "./assets"),
	}
}

// Scaffold a koa app
func (d *D) Scaffold(output string) error {
	return bundler.Restore(d.assets, output)
}

// Bundle a function into a koa project
func (d *D) Bundle(output string, fn string, deps ...string) error {
	return bundler.Bundle(d.assets, output, "d", fn, deps...)
}

var (
	_ bundler.Bundler = &D{}
)
