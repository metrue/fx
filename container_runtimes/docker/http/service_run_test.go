package api

import (
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/metrue/fx/config"
	"github.com/metrue/fx/constants"
	"github.com/metrue/fx/types"

	gock "gopkg.in/h2non/gock.v1"
)

func TestServiceRun(t *testing.T) {
	defer gock.Off()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	host := config.Host{Host: "127.0.0.1"}
	api, err := Create(host.Host, constants.AgentPort)
	if err != nil {
		t.Fatal(err)
	}

	service := types.Service{
		Name:  "a-mock-service",
		Image: "a-mock-image-id",
	}

	mockContainerID := "mock-container-id"
	url := "http://" + host.Host + ":" + constants.AgentPort
	gock.New(url).
		Post("/v0.2.1/containers").
		AddMatcher(func(req *http.Request, ereq *gock.Request) (m bool, e error) {
			// TODO multiple matching not supported by gock
			if req.URL.String() == url+"/v0.2.1/containers/"+mockContainerID+"/start" {
				return true, nil
			} else if req.URL.String() == url+"/v0.2.1/containers/create?name="+service.Name {
				return true, nil
			}

			return false, nil
		}).
		Reply(201).
		JSON(map[string]interface{}{
			"Id":       mockContainerID,
			"Warnings": []string{},
		})

	// FIXME
	if err := api.Run(9999, &service); err == nil {
		t.Fatal(err)
	}
}
