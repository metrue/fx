package context

import (
	"testing"

	"github.com/urfave/cli"
)

func TestContext(t *testing.T) {
	ctx := NewContext()
	cli := cli.NewContext(nil, nil, nil)
	ctx.WithCliContext(cli)
	c := ctx.GetCliContext()
	if c != cli {
		t.Fatalf("should get %v but got %v", cli, c)
	}

	key := "k_1"
	value := "hello"
	ctx.Set(key, "hello")
	v := ctx.Get(key).(string)
	if v != value {
		t.Fatalf("should get %v but %v", value, v)
	}
}
