package middlewares

import (
	"os"
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
		pwd, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}
		argset.Parse([]string{pwd})
		ctx.EXPECT().GetCliContext().Return(cli)
		ctx.EXPECT().Set("sources", []string{pwd})
		ctx.EXPECT().Set("name", "")
		ctx.EXPECT().Set("port", 0)
		ctx.EXPECT().Set("force", false)
		if err := Parse("up")(ctx); err != nil {
			t.Fatal("should got file or directory not existed error")
		}
	})
}
