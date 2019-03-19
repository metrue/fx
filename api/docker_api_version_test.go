package api

import (
	"testing"

	gock "gopkg.in/h2non/gock.v1"
)

func TestDockerAPIVersion(t *testing.T) {
	defer gock.Off()

	dockerRemoteAPI := "http://127.0.0.1:1234"
	gock.New(dockerRemoteAPI).
		Get("/version").
		Reply(200).
		JSON(map[string]string{
			"ApiVersion": "0.2.1",
		})
	v, err := Version(dockerRemoteAPI)
	if err != nil {
		t.Fatal(err)
	}
	if v != "0.2.1" {
		t.Fatalf("should get %s but got %s", "0.2.1", v)
	}
}
