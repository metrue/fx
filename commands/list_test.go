package commands_test

import (
	"testing"

	. "github.com/metrue/fx/commands"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	addr := "localhost:5000"
	functions := []string{"*"}

	err := List(addr, functions)
	assert.NotNil(t, err)
}
