package handlers

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	mockCtx "github.com/metrue/fx/context/mocks"
	mockDeployer "github.com/metrue/fx/infra/mocks"
	"github.com/metrue/fx/types"
)

func TestUp(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := mockCtx.NewMockContexter(ctrl)
	deployer := mockDeployer.NewMockDeployer(ctrl)

	bindings := []types.PortBinding{}
	name := "sample-name"
	image := "sample-image"
	data := "sample-data"
	ctx.EXPECT().Get("name").Return(name)
	ctx.EXPECT().Get("image").Return(image)
	ctx.EXPECT().Get("deployer").Return(deployer)
	ctx.EXPECT().Get("bindings").Return(bindings)
	ctx.EXPECT().Get("data").Return(data)
	ctx.EXPECT().GetContext().Return(context.Background()).Times(2)
	deployer.EXPECT().Deploy(gomock.Any(), data, name, image, bindings).Return(nil)
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
