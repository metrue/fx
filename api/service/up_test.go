package service_test

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/metrue/fx/api"
	. "github.com/metrue/fx/api/service"
	"github.com/stretchr/testify/assert"
)

func TestUp(t *testing.T) {
	meta := api.FunctionMeta{
		Lang: "node",
		Path: "./func.js",
		Content: `module.exports = function (input) {
	return input.a + input.b
}
		`,
	}
	ret, err := DoUp(meta)
	assert.Nil(t, nil, err)

	time.Sleep(5 * time.Second)

	body := strings.NewReader(`{"a": 1, "b": 1}`)

	u, err := url.Parse(fmt.Sprintf("http://%s", ret.LocalAddress))
	if err != nil {
		log.Fatal(err)
	}

	resp, err := http.Post(u.String(), "application/json", body)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	assert.Equal(t, buf.String(), "2")
}
