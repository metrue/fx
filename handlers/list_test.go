package handlers_test

import (
	"reflect"
	"testing"

	. "github.com/metrue/fx/handlers"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	ids := ""
	contains, err := List(ids)
	assert.Nil(t, err)
	assert.Equal(t, reflect.TypeOf(contains).Kind(), reflect.Slice)
}
