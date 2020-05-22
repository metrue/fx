package rust

import (
	"github.com/gobuffalo/packr/v2"
	"github.com/metrue/fx/bundler"
)

// Rust defines javascript bundler
type Rust struct {
	assets *packr.Box
}

// New a koa bundler
func New() *Rust {
	return &Rust{
		assets: packr.New("rust", "./assets"),
	}
}

// Scaffold a koa app
func (k *Rust) Scaffold(output string) error {
	return bundler.Restore(k.assets, output)
}

// Bundle a function into a koa project
func (k *Rust) Bundle(output string, fn string, deps ...string) error {
	return bundler.Bundle(k.assets, output, "rust", fn, deps...)
}

var (
	_ bundler.Bundler = &Rust{}
)
