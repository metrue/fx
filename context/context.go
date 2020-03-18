package context

import (
	"context"

	"github.com/urfave/cli"
)

type key string

const (
	keyCliCtx = key("cmd_cli")
)

// Contexter ctx interface
type Contexter interface {
	Get(k string) interface{}
	Set(k string, v interface{})
	Use(fn func(ctx *Context) error) error
	GetContext() context.Context
	GetCliContext() *cli.Context
}

// Context fx context
type Context struct {
	context.Context
}

// NewContext new a context
func NewContext() *Context {
	ctx := context.Background()
	return &Context{ctx}
}

// FromCliContext create context from cli.Context
func FromCliContext(c *cli.Context) *Context {
	ctx := NewContext()
	ctx.WithCliContext(c)
	return ctx
}

// WithCliContext set cli.Context
func (ctx *Context) WithCliContext(c *cli.Context) {
	newCtx := context.WithValue(ctx.Context, keyCliCtx, c)
	ctx.Context = newCtx
}

// GetCliContext get cli.Context
func (ctx *Context) GetCliContext() *cli.Context {
	return ctx.Value(keyCliCtx).(*cli.Context)
}

// Set a value with name
func (ctx *Context) Set(name string, value interface{}) {
	// nolint
	newCtx := context.WithValue(ctx.Context, name, value)
	ctx.Context = newCtx
}

// Get a value
func (ctx *Context) Get(name string) interface{} {
	return ctx.Context.Value(name)
}

// Use invole a middle
func (ctx *Context) Use(fn func(ctx *Context) error) error {
	return fn(ctx)
}

// GetContext get context
func (ctx *Context) GetContext() context.Context {
	return ctx.Context
}

var (
	_ Contexter = &Context{}
)
