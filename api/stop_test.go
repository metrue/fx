package api

import (
	"testing"

	gock "gopkg.in/h2non/gock.v1"
)

func TestStop(t *testing.T) {
	defer gock.Off()

	dockerRemoteAPI := "http://127.0.0.1:1234"
	version := "0.2.1"
	api := NewWithDockerRemoteAPI(dockerRemoteAPI, version)

	mockServiceName := "mock-service-name"
	gock.New(dockerRemoteAPI).
		Post("/v0.2.1/containers/" + mockServiceName + "/stop").
		Reply(204)
	if err := api.Stop(mockServiceName); err != nil {
		t.Fatal(err)
	}
}
