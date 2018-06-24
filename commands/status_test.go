package commands_test

import (
	"testing"
	"time"

	. "github.com/metrue/fx/commands"
	"github.com/metrue/fx/pkg/server"
	"github.com/stretchr/testify/assert"
)

func TestStatus(t *testing.T) {
	addr := "localhost:23453"
	s := server.NewFxServiceServer(addr)
	go func() {
		s.Start()
	}()
	time.Sleep(2 * time.Second)

	err := Status(addr)
	assert.Nil(t, err)

	s.Stop()
	time.Sleep(2 * time.Second)
}
