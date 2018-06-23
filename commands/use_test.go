package commands_test

import (
	"testing"

	. "github.com/metrue/fx/commands"
	"github.com/metrue/fx/config"

	"github.com/stretchr/testify/assert"
)

func TestUse(t *testing.T) {
	config.CONFIG = "/tmp/fx.json"

	addr := "a.b.c.d:124"
	err := Use(addr)
	assert.Nil(t, err)

	newConf := config.GetConfig()
	assert.Equal(t, newConf.HttpServerAddr, "a.b.c.d:124")
	assert.Equal(t, newConf.GrpcEndpoint, "a.b.c.d:124")
}
