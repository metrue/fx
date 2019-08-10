package api

import (
	"testing"

	"github.com/gobuffalo/packr"
	"github.com/golang/mock/gomock"
	"github.com/metrue/fx/config"
	"github.com/metrue/fx/constants"
	gock "gopkg.in/h2non/gock.v1"
)

func TestDockerAPIVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	defer gock.Off()

	host := config.Host{Host: "127.0.0.1"}
	const version = "0.2.1"

	box := packr.NewBox("./api/images")
	api := New(box)
	if err := api.Init(host); err != nil {
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
