package up

import (
	"log"
	"os"
)

// Customer logger
type Logger struct {
	base *log.Logger
}

func (l *Logger) Log(v interface{}) {
	l.base.SetOutput(os.Stdout)
	l.base.Print(v)
}

func (l *Logger) Err(v interface{}) {
	l.base.SetOutput(os.Stderr)
	l.base.Print(v)
}

func NewLogger(prefix string) *Logger {
	return &Logger{
		base: log.New(
			os.Stdout,
			prefix,
			log.LstdFlags,
		),
	}
}
