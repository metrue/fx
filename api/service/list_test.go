package service_test

import (
	"reflect"
	"testing"

	. "github.com/metrue/fx/api/service"
	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	ids := ""
	contains, err := DoList(ids)
	assert.Nil(t, err)
	assert.Equal(t, reflect.TypeOf(contains).Kind(), reflect.Slice)
}
