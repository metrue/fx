package middlewares

import (
	"testing"

	"github.com/golang/mock/gomock"
	mockCtx "github.com/metrue/fx/context/mocks"
)

func TestProvision(t *testing.T) {
	t.Skip()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := mockCtx.NewMockContexter(ctrl)

	ctx.EXPECT().Get("host").Return("root@127.0.0.1")
	ctx.EXPECT().Get("kubeconf").Return("~/.kube/config")
	ctx.EXPECT().Get("ssh_port").Return("22")
	ctx.EXPECT().Get("ssh_key").Return("~/.ssh/id_rsa")

	if err := Provision(ctx); err != nil {
		t.Fatal(err)
	}
}
