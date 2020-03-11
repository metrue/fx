package middlewares

import (
	"testing"

	"github.com/golang/mock/gomock"
	mockCtx "github.com/metrue/fx/context/mocks"
)

func TestLanguage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := mockCtx.NewMockContexter(ctrl)
	ctx.EXPECT().Get("fn").Return("/tmp/fx.js")
	ctx.EXPECT().Set("language", "node")
	if err := Language()(ctx); err != nil {
		t.Fatal(err)
	}
}
