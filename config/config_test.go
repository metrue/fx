package config_test

import (
	"os"
	"testing"

	. "github.com/metrue/fx/config"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {

	c := GetConfig()
	assert.NotNil(t, c)

	assert.Equal(t, GetHttpServerAddr(), "localhost:30080")
	assert.Equal(t, GetGrpcEndpoint(), "localhost:50000")
}

func TestSetHost(t *testing.T) {
	CONFIG = "/tmp/fx.config.json"

	c := GetConfig()
	assert.Equal(t, GetHttpServerAddr(), "localhost:30080")
	assert.Equal(t, GetGrpcEndpoint(), "localhost:50000")

	err := c.SetHost("124.124.124.124")
	assert.Nil(t, err)
	assert.Equal(t, "124.124.124.124:30080", GetHttpServerAddr())
	assert.Equal(t, "124.124.124.124:50000", GetGrpcEndpoint())
}

func cleanup() {
	CONFIG = "/tmp/fx.config.json"
	os.Remove(CONFIG)
}

func TestMain(m *testing.M) {
	cleanup()

	m.Run()

	// cleanup()
}
