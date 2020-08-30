package crystal

import (
	"github.com/gobuffalo/packr/v2"
	"github.com/metrue/fx/bundler"
)

// Crystal defines crystal bundler
type Crystal struct {
	assets *packr.Box
}

// New a crystal bundler
func New() *Crystal {
	return &Crystal{
		assets: packr.New("crystal", "./assets"),
	}
}

// Scaffold a crystal/kemal app
func (k *Crystal) Scaffold(output string) error {
	return bundler.Restore(k.assets, output)
}

// Bundle a function into a kemal project
func (k *Crystal) Bundle(output string, fn string, deps ...string) error {
	return bundler.Bundle(k.assets, output, "crystal", fn, deps...)
}

var (
	_ bundler.Bundler = &Crystal{}
)
