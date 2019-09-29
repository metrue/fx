package rkt

import (
	"github.com/metrue/fx/builders"
)

// Rkt rkt as a image builder
type Rkt struct{}

// Build build a directory to be a image
func (r *Rkt) Build(workdir string, name string) error {
	return nil
}

var (
	_ builders.Builder = &Rkt{}
)
