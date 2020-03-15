package julia

import (
	"github.com/gobuffalo/packr/v2"
	"github.com/metrue/fx/bundler"
)

// Julia defines javascript bundler
type Julia struct {
	assets *packr.Box
}

// New a koa bundler
func New() *Julia {
	return &Julia{
		assets: packr.New("julia", "./assets"),
	}
}

// Scaffold a koa app
func (k *Julia) Scaffold(output string) error {
	return bundler.Restore(k.assets, output)
}

// Bundle a function into a koa project
func (k *Julia) Bundle(output string, fn string, deps ...string) error {
	return bundler.Bundle(k.assets, output, "julia", fn, deps...)
}

var (
	_ bundler.Bundler = &Julia{}
)
