package commands_test

import (
	"testing"
	"time"

	. "github.com/metrue/fx/commands"
	"github.com/metrue/fx/pkg/server"
	"github.com/stretchr/testify/assert"
)

func TestDown(t *testing.T) {
	addr := "localhost:23451"
	functions := []string{"id-should-not-exist"}

	s := server.NewFxServiceServer(addr)
	go func() {
		s.Start()
	}()
	time.Sleep(2 * time.Second)

	err := Down(addr, functions)
	assert.Nil(t, err)

	s.Stop()
	time.Sleep(2 * time.Second)
}
