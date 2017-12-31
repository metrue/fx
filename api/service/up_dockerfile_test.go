package service

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/metrue/fx/api"
)

func createDockerfile(client api.FxServiceClient) (*api.UpMsgMeta, error) {
	ctx := context.Background()
	upReq := &api.UpDockerfileRequest{
		Dockerfiles: []*api.DockerfileMeta{
			&api.DockerfileMeta{
				Content: fmt.Sprintf(`FROM andygrunwald/simple-webserver
RUN echo '%d'`, time.Now().Unix()),
			},
		},
	}

	fmt.Println("to Up 1")
	upRes, err := client.UpDockerfile(ctx, upReq)
	if err != nil {
		return nil, err
	}
	fmt.Println("to Up 2")

	if len(upRes.Instances) != 1 {
		return nil, fmt.Errorf("UpDockerfile response should have one instance, found %d", len(upRes.Instances))
	}
	fmt.Println("to Up 3")
	if upRes.Instances[0].Error != "" {
		return nil, fmt.Errorf("Up error: %s", upRes.Instances[0].Error)
	}

	fmt.Println("to Up 4")
	return upRes.Instances[0], nil
}

func TestUpDockerfile(t *testing.T) {
	runServer(t)

	client, conn, err := api.NewClient(grpcEndpoint)
	defer stopServer(conn)

	if err != nil {
		t.Fatal(err)
	}

	_, createErr := createDockerfile(client)
	if createErr != nil {
		t.Fatal("Up Dockerfile: %s\n", err.Error())
	}
}
