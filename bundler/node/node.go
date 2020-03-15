package node

import (
	"github.com/gobuffalo/packr/v2"
	"github.com/metrue/fx/bundler"
)

const language = "node"

// Node defines node bundler
type Node struct {
	assets *packr.Box
}

// New a koa bundler
func New() *Node {
	return &Node{
		assets: packr.New("node", "./assets"),
	}
}

// Scaffold a koa app
func (k *Node) Scaffold(output string) error {
	return bundler.Restore(k.assets, output)
}

// Bundle a function into a koa project
func (k *Node) Bundle(output string, fn string, deps ...string) error {
	return bundler.Bundle(k.assets, output, language, fn, deps...)
}

var (
	_ bundler.Bundler = &Node{}
)
