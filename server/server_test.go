package server_test

import (
	"testing"

	. "github.com/metrue/fx/server"
	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	err := Start(true)
	assert.Nil(t, err)
}

func TestMain(m *testing.M) {
	m.Run()
}
