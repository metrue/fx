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
	t.Run("normally up", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := mockCtx.NewMockContexter(ctrl)
		driver := mockDeployer.NewMockDeployer(ctrl)

		bindings := []types.PortBinding{}
		name := "sample-name"
		image := "sample-image"
		data := "sample-data"
		ctx.EXPECT().Get("name").Return(name)
		ctx.EXPECT().Get("image").Return(image)
		ctx.EXPECT().Get("docker_driver").Return(driver)
		ctx.EXPECT().Get("k8s_driver").Return(driver)
		ctx.EXPECT().Get("bindings").Return(bindings)
		ctx.EXPECT().Get("data").Return(data)
		ctx.EXPECT().Get("force").Return(false)
		ctx.EXPECT().GetContext().Return(context.Background()).Times(4)
		driver.EXPECT().Deploy(gomock.Any(), data, name, image, bindings).Return(nil).Times(2)
		driver.EXPECT().GetStatus(gomock.Any(), name).Return(types.Service{
			ID:   "id-1",
			Name: name,
			Host: "127.0.0.1",
			Port: 2100,
		}, nil).Times(2)
		if err := Up(ctx); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("normally up forcely", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := mockCtx.NewMockContexter(ctrl)
		driver := mockDeployer.NewMockDeployer(ctrl)

		bindings := []types.PortBinding{}
		name := "sample-name"
		image := "sample-image"
		data := "sample-data"
		ctx.EXPECT().Get("name").Return(name)
		ctx.EXPECT().Get("image").Return(image)
		ctx.EXPECT().Get("docker_driver").Return(driver)
		ctx.EXPECT().Get("k8s_driver").Return(driver)
		ctx.EXPECT().Get("bindings").Return(bindings)
		ctx.EXPECT().Get("data").Return(data)
		ctx.EXPECT().Get("force").Return(true)
		ctx.EXPECT().GetContext().Return(context.Background()).Times(6)
		driver.EXPECT().Deploy(gomock.Any(), data, name, image, bindings).Return(nil).Times(2)
		driver.EXPECT().Destroy(gomock.Any(), name).Return(nil).Times(2)
		driver.EXPECT().GetStatus(gomock.Any(), name).Return(types.Service{
			ID:   "id-1",
			Name: name,
			Host: "127.0.0.1",
			Port: 2100,
		}, nil).Times(2)
		if err := Up(ctx); err != nil {
			t.Fatal(err)
		}
	})
}
