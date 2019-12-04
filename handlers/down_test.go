package handlers

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	mockCtx "github.com/metrue/fx/context/mocks"
	mockDeployer "github.com/metrue/fx/infra/mocks"
)

func TestDown(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := mockCtx.NewMockContexter(ctrl)
	deployer := mockDeployer.NewMockDeployer(ctrl)

	services := []string{"sample-name"}
	ctx.EXPECT().Get("services").Return(services)
	ctx.EXPECT().Get("deployer").Return(deployer)
	ctx.EXPECT().GetContext().Return(context.Background())
	deployer.EXPECT().Destroy(gomock.Any(), services[0]).Return(nil)
	if err := Down(ctx); err != nil {
		t.Fatal(err)
	}
}
