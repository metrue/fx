package php

import (
	"github.com/gobuffalo/packr/v2"
	"github.com/metrue/fx/bundler"
)

// Php defines php bundler
type Php struct {
	assets *packr.Box
}

// New a koa bundler
func New() *Php {
	return &Php{
		assets: packr.New("php", "./assets"),
	}
}

// Scaffold a koa app
func (k *Php) Scaffold(output string) error {
	return bundler.Restore(k.assets, output)
}

// Bundle a function into a koa project
func (k *Php) Bundle(output string, fn string, deps ...string) error {
	return bundler.Bundle(k.assets, output, "php", fn, deps...)
}

var (
	_ bundler.Bundler = &Php{}
)
