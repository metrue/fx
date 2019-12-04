package handlers

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	mockCtx "github.com/metrue/fx/context/mocks"
	mockDeployer "github.com/metrue/fx/infra/mocks"
	"github.com/metrue/fx/types"
	fxTypes "github.com/metrue/fx/types"
)

func TestUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := mockCtx.NewMockContexter(ctrl)
	deployer := mockDeployer.NewMockDeployer(ctrl)

	fn := fxTypes.Func{}
	bindings := []types.PortBinding{}
	name := "sample-name"
	image := "sample-image"
	ctx.EXPECT().Get("fn").Return(fn)
	ctx.EXPECT().Get("name").Return(name)
	ctx.EXPECT().Get("image").Return(image)
	ctx.EXPECT().Get("deployer").Return(deployer)
	ctx.EXPECT().Get("bindings").Return(bindings)
	ctx.EXPECT().GetContext().Return(context.Background()).Times(2)
	deployer.EXPECT().Deploy(gomock.Any(), fn, name, image, bindings).Return(nil)
	deployer.EXPECT().GetStatus(gomock.Any(), name).Return(types.Service{
		ID:   "id-1",
		Name: name,
		Host: "127.0.0.1",
		Port: 2100,
	}, nil)
	if err := Up(ctx); err != nil {
		t.Fatal(err)
	}
}
