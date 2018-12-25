package commands_test

import (
	"testing"
	"time"

	. "github.com/metrue/fx/commands"
	"github.com/metrue/fx/pkg/server"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	addr := "localhost:23453"
	s := server.NewFxServiceServer(addr)
	go func() {
		s.Start()
	}()
	time.Sleep(2 * time.Second)

	functions := []string{"*"}

	err := List(addr, functions)
	assert.Nil(t, err)

	s.Stop()
	time.Sleep(2 * time.Second)
}
