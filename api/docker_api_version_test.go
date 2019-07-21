package api

import (
	"testing"

	"github.com/metrue/fx/constants"
	gock "gopkg.in/h2non/gock.v1"
)

func TestDockerAPIVersion(t *testing.T) {
	defer gock.Off()

	const version = "0.2.1"

	dockerRemoteAPI := "http://" + constants.DockerRemoteAPIEndpoint
	gock.New(dockerRemoteAPI).
		Get("/version").
		Reply(200).
		JSON(map[string]string{
			"ApiVersion": version,
		})
	v, err := Version(dockerRemoteAPI)
	if err != nil {
		t.Fatal(err)
	}
	if v != version {
		t.Fatalf("should get %s but got %s", version, v)
	}
}
