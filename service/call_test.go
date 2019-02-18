package service_test

import (
	"testing"

	"github.com/metrue/fx/api"
	. "github.com/metrue/fx/service"
	"github.com/stretchr/testify/assert"
)

func TestCall(t *testing.T) {
	req := &api.CallRequest{
		Lang: "node",
		Content: `module.exports = function (input) {
	return parseInt(input.a, 10) + parseInt(input.b, 10)
}
		`,
		Params: "a=1 b=4",
	}
	res, err := Call(nil, req)
	assert.Nil(t, err)
	assert.Equal(t, "5", res.Data)
}
