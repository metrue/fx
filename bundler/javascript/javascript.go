package javascript

import (
	"github.com/gobuffalo/packr/v2"
	"github.com/metrue/fx/bundler"
)

const language = "node"

// Koa defines javascript bundler
type Koa struct {
	assets *packr.Box
}

// New a koa bundler
func New() *Koa {
	return &Koa{
		assets: packr.New("", "./assets"),
	}
}

// Scaffold a koa app
func (k *Koa) Scaffold(output string) error {
	return bundler.Restore(k.assets, output)
}

// Bundle a function into a koa project
func (k *Koa) Bundle(output string, fn string, deps ...string) error {
	return bundler.Bundle(k.assets, output, language, fn, deps...)
}

var (
	_ bundler.Bundler = &Koa{}
)
