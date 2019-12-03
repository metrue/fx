package middlewares

import (
	"testing"

	"github.com/golang/mock/gomock"
	mockCtx "github.com/metrue/fx/context/mocks"
)

func TestBinding(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := mockCtx.NewMockContexter(ctrl)
	ctx.EXPECT().Get("port").Return(0)
	ctx.EXPECT().Set("bindings", gomock.Any())
	if err := Binding(ctx); err != nil {
		t.Fatal(err)
	}
}
