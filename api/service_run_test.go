package api

import (
	"net/http"
	"testing"

	"github.com/metrue/fx/types"

	gock "gopkg.in/h2non/gock.v1"
)

func TestServiceRun(t *testing.T) {
	defer gock.Off()

	dockerRemoteAPI := "http://127.0.0.1:1234"
	version := "0.2.1"

	service := types.Service{
		Name:  "a-mock-service",
		Image: "a-mock-image-id",
	}

	mockContainerID := "mock-container-id"
	gock.New(dockerRemoteAPI).
		Post("/v0.2.1/containers").
		AddMatcher(func(req *http.Request, ereq *gock.Request) (m bool, e error) {
			// TODO multiple matching not supported by gock
			if req.URL.String() == dockerRemoteAPI+"/v0.2.1/containers/"+mockContainerID+"/start" {
				return true, nil
			} else if req.URL.String() == dockerRemoteAPI+"/v0.2.1/containers/create?name="+service.Name {
				return true, nil
			}

			return false, nil
		}).
		Reply(201).
		JSON(map[string]interface{}{
			"Id":       mockContainerID,
			"Warnings": []string{},
		})

	api := NewWithDockerRemoteAPI(dockerRemoteAPI, version)
	// FIXME
	if err := api.Run(&service); err == nil {
		t.Fatal(err)
	}
}
