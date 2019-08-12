package handlers

import "github.com/urfave/cli"

// HandleFunc command handle function
type HandleFunc func(ctx *cli.Context) error
