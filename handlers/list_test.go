package handlers_test

import (
	"testing"

	. "github.com/metrue/fx/handlers"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	ids := ""
	contains, err := List(ids)
	assert.Nil(t, err)
}
