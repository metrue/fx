package api

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gobuffalo/packr"
	"github.com/golang/mock/gomock"
	"github.com/metrue/fx/config"
	mockConfig "github.com/metrue/fx/config/mocks"
	"github.com/metrue/fx/constants"
	gock "gopkg.in/h2non/gock.v1"
)

func TestStop(t *testing.T) {
	defer gock.Off()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	host := config.Host{Host: "127.0.0.1"}
	cfg := mockConfig.NewMockConfiger(ctrl)
	cfg.EXPECT().GetDefaultHost().Return(host, nil)
	box := packr.NewBox("./images")
	api := New(cfg, box)
	if err := api.Init(); err != nil {
		t.Fatal(err)
	}

	mockServiceName := "mock-service-name"
	url := "http://" + host.Host + ":" + constants.AgentPort
	gock.New(url).
		Post("/v" + api.version + "/containers/" + mockServiceName + "/stop").
		AddMatcher(func(req *http.Request, ereq *gock.Request) (m bool, e error) {
			if strings.Contains(req.URL.String(), "/v"+api.version+"/containers/"+mockServiceName+"/stop") {
				return true, nil
			}
			return false, nil
		}).
		Reply(204)
	if err := api.Stop(mockServiceName); err != nil {
		t.Fatal(err)
	}
}
