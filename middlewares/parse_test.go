package middlewares

import (
	"io/ioutil"
	"os"
	"os/user"
	"testing"

	"flag"

	"github.com/golang/mock/gomock"
	mockCtx "github.com/metrue/fx/context/mocks"
	"github.com/urfave/cli"
)

func TestParse(t *testing.T) {
	t.Run("source code not existed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := mockCtx.NewMockContexter(ctrl)
		argset := flag.NewFlagSet("test", 0)
		cli := cli.NewContext(nil, argset, nil)
		argset.Parse([]string{"this_file_should_not_existed"})
		ctx.EXPECT().GetCliContext().Return(cli)
		if err := Parse("up")(ctx); err == nil {
			t.Fatal("should got file or directory not existed error")
		}
	})
	t.Run("source code ready", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := mockCtx.NewMockContexter(ctrl)
		argset := flag.NewFlagSet("test", 0)
		cli := cli.NewContext(nil, argset, nil)
		fd, err := ioutil.TempFile("", "fx_func_*.js")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(fd.Name())

		user, err := user.Current()
		if err != nil {
			t.Fatal(err)
		}
		host := user.Username + "@localhost"

		argset.Parse([]string{fd.Name()})
		ctx.EXPECT().GetCliContext().Return(cli)
		ctx.EXPECT().Set("fn", fd.Name())
		ctx.EXPECT().Set("deps", []string{})
		ctx.EXPECT().Set("host", host)
		ctx.EXPECT().Set("kubeconf", "")
		ctx.EXPECT().Set("name", "")
		ctx.EXPECT().Set("port", 0)
		ctx.EXPECT().Set("force", false)
		if err := Parse("up")(ctx); err != nil {
			t.Fatal(err)
		}
	})
}
