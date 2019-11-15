package handlers

import (
	"github.com/metrue/fx/context"
)

// HandleFunc command handle function
type HandleFunc func(ctx *context.Context) error
