package api

import (
	"testing"

	"github.com/gobuffalo/packr"
	"github.com/golang/mock/gomock"
	"github.com/metrue/fx/config"
	mockConfig "github.com/metrue/fx/config/mocks"
	"github.com/metrue/fx/constants"
	gock "gopkg.in/h2non/gock.v1"
)

func TestDockerAPIVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	defer gock.Off()

	host := config.Host{Host: "127.0.0.1"}
	const version = "0.2.1"

	cfg := mockConfig.NewMockConfiger(ctrl)
	cfg.EXPECT().GetDefaultHost().Return(host, nil)
	box := packr.NewBox("./api/images")
	api := New(cfg, box)
	if err := api.Init(); err != nil {
		t.Fatal(err)
	}

	url := "http://" + host.Host + ":" + constants.AgentPort
	gock.New(url).
		Get("/version").
		Reply(200).
		JSON(map[string]string{
			"ApiVersion": version,
		})
	v, err := api.Version(url)
	if err != nil {
		t.Fatal(err)
	}
	if v != version {
		t.Fatalf("should get %s but got %s", version, v)
	}
}
