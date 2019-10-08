package api

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/metrue/fx/config"
	"github.com/metrue/fx/constants"
	gock "gopkg.in/h2non/gock.v1"
)

func TestNetwork(t *testing.T) {
	defer gock.Off()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	host := config.Host{Host: "127.0.0.1"}
	api, err := Create(host.Host, constants.AgentPort)
	if err != nil {
		t.Fatal(err)
	}

	const network = "fx-net"
	if err := api.CreateNetwork(network); err != nil {
		t.Fatal(err)
	}

	nws, err := api.GetNetwork(network)
	if err != nil {
		t.Fatal(err)
	}
	if nws[0].Name != network {
		t.Fatalf("should get %s but got %s", network, nws[0].Name)
	}
}
