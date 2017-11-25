package log

import (
	"log"
	"os"
)

// Logger is a customer wrapped logger
type Logger struct {
	base *log.Logger
	log.Logger
}

// Log outputs an object to stdout
func (l *Logger) Log(v interface{}) {
	l.base.SetOutput(os.Stdout)
	l.base.Print(v)
}

// Err outputs an object to stderr
func (l *Logger) Err(v interface{}) {
	l.base.SetOutput(os.Stderr)
	l.base.Print(v)
}

// NewLogger creates a new wrapped log.Logger with default values
func NewLogger(prefix string) *Logger {
	return &Logger{
		base: log.New(
			os.Stdout,
			prefix,
			log.LstdFlags,
		),
	}
}
