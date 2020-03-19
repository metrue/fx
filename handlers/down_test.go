package handlers

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	mockCtx "github.com/metrue/fx/context/mocks"
	mockDeployer "github.com/metrue/fx/driver/mocks"
)

func TestDown(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := mockCtx.NewMockContexter(ctrl)
	driver := mockDeployer.NewMockDriver(ctrl)

	services := []string{"sample-name"}
	ctx.EXPECT().Get("services").Return(services)
	ctx.EXPECT().Get("docker_driver").Return(driver)
	ctx.EXPECT().Get("k8s_driver").Return(driver)
	ctx.EXPECT().GetContext().Return(context.Background()).Times(2)
	driver.EXPECT().Destroy(gomock.Any(), services[0]).Return(nil).Times(2)
	if err := Down(ctx); err != nil {
		t.Fatal(err)
	}
}
