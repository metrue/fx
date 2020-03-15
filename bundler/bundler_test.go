package bundler

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/gobuffalo/packr/v2"
)

func TestBundler(t *testing.T) {
	t.Run("Restore", func(t *testing.T) {
		box := packr.New("", "./node/assets")
		outputDir, err := ioutil.TempDir("", "fx_koa")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(outputDir)
		if err := Restore(box, outputDir); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Bundle", func(t *testing.T) {
		fd, err := ioutil.TempFile("", "fx_func_*.js")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(fd.Name())

		content := `
module.exports = (ctx) => {
  ctx.body = 'hello fx'
}`
		if err = ioutil.WriteFile(fd.Name(), []byte(content), 0666); err != nil {
			t.Fatal(err)
		}

		outputDir, err := ioutil.TempDir("", "fx_koa")
		if err != nil {
			log.Fatal(err)
		}
		defer os.RemoveAll(outputDir)

		box := packr.New("", "./node/assets")
		if err := Bundle(box, outputDir, "node", fd.Name()); err != nil {
			t.Fatal(err)
		}
	})
}
