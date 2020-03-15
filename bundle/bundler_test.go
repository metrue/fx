package bundle

import (
	"io/ioutil"
	"os"
	"testing"
)

func createFn(content string, t *testing.T) string {
	fd, err := ioutil.TempFile("", "fx_func_*.js")
	if err != nil {
		t.Fatal(err)
	}
	err = ioutil.WriteFile(fd.Name(), []byte(content), 0666)
	if err != nil {
		t.Fatal(err)
	}

	return fd.Name()
}

func TestBundle(t *testing.T) {
	workdir, err := ioutil.TempDir("", "fx-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(workdir)

	cases := []struct {
		workdir  string
		language string
		fn       string
		deps     []string
	}{
		{
			workdir:  workdir,
			language: "d",
			fn: `
module.exports = (ctx) => {
  ctx.body = 'hello fx'
}
			`,
		},
		{
			workdir:  workdir,
			language: "go",
			fn: `
module.exports = (ctx) => {
  ctx.body = 'hello fx'
}
			`,
		},
		{
			workdir:  workdir,
			language: "java",
			fn: `
module.exports = (ctx) => {
  ctx.body = 'hello fx'
}
			`,
		},
		{
			workdir:  workdir,
			language: "julia",
			fn: `
module.exports = (ctx) => {
  ctx.body = 'hello fx'
}
			`,
		},
		{
			workdir:  workdir,
			language: "perl",
			fn: `
module.exports = (ctx) => {
  ctx.body = 'hello fx'
}
			`,
		},
		{
			workdir:  workdir,
			language: "python",
			fn: `
module.exports = (ctx) => {
  ctx.body = 'hello fx'
}
			`,
		},
		{
			workdir:  workdir,
			language: "ruby",
			fn: `
module.exports = (ctx) => {
  ctx.body = 'hello fx'
}
			`,
		},
		{
			workdir:  workdir,
			language: "rust",
			fn: `
module.exports = (ctx) => {
  ctx.body = 'hello fx'
}
			`,
		},
	}

	for _, c := range cases {
		fn := createFn(c.fn, t)
		defer os.Remove(fn)

		if err := Bundle(c.workdir, c.language, fn, c.deps...); err != nil {
			t.Fatal(err)
		}
	}
}
