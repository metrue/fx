package commands_test

import (
	"testing"

	. "github.com/metrue/fx/commands"
	"github.com/stretchr/testify/assert"
)

func TestDown(t *testing.T) {
	addr := "localhost:30080"
	functions := []string{"id-should-not-exist"}

	err := Down(addr, functions)
	assert.Equal(t, DownFunctionError, err)
}
