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
	assert.Equal(t, c.HttpServerAddr, "localhost:30080")
	assert.Equal(t, c.GrpcEndpoint, "localhost:5000")
}

func TestSave(t *testing.T) {
	CONFIG = "/tmp/fx.config.json"

	c := GetConfig()
	assert.Equal(t, c.HttpServerAddr, "localhost:30080")
	assert.Equal(t, c.GrpcEndpoint, "localhost:5000")

	c.HttpServerAddr = "123.123.123.123:1234"
	c.GrpcEndpoint = "321.321.321.321:4321"

	err := c.Save()
	assert.Nil(t, err)

	newConfig := GetConfig()

	assert.Equal(t, newConfig.HttpServerAddr, "123.123.123.123:1234")
	assert.Equal(t, newConfig.GrpcEndpoint, "321.321.321.321:4321")
}

func cleanup() {
	CONFIG = "/tmp/fx.config.json"
	os.Remove(CONFIG)
}

func TestMain(m *testing.M) {
	cleanup()

	m.Run()

	cleanup()
}
