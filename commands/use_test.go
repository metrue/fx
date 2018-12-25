package commands_test

import (
	"testing"

	. "github.com/metrue/fx/commands"
	"github.com/metrue/fx/config"

	"github.com/stretchr/testify/assert"
)

func TestUse(t *testing.T) {
	config.CONFIG = "/tmp/fx.json"

	addr := "a.b.c.d"
	err := Use(addr)
	assert.Nil(t, err)

	assert.Equal(t, config.GetHttpServerAddr(), "a.b.c.d:30080")
	assert.Equal(t, config.GetGrpcEndpoint(), "a.b.c.d:50000")
}
