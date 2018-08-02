package commands_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/metrue/fx/commands"
	"github.com/metrue/fx/pkg/server"
	"github.com/stretchr/testify/assert"
)

func TestCall(t *testing.T) {
	addr := "localhost:23451"

	s := server.NewFxServiceServer(addr)
	go func() {
		s.Start()
	}()
	time.Sleep(2 * time.Second)
	defer s.Stop()

	code := `
module.exports = (input) => {
		console.log("acc")
    return parseInt(input.a, 10) + parseInt(input.b, 10)
}
`
	content := []byte(code)
	tmpDir, err := ioutil.TempDir("", "fx")
	assert.Nil(t, err)
	assert.NotEqual(t, "", tmpDir)

	defer os.RemoveAll(tmpDir)

	tmpfn := filepath.Join(tmpDir, "func.js")
	err = ioutil.WriteFile(tmpfn, content, 0666)
	assert.Nil(t, err)

	params := "a=1 b=1"
	err = Call(addr, tmpfn, params)
	assert.Nil(t, err)

	time.Sleep(2 * time.Second)
}
