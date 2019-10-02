package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/apex/log"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/metrue/fx/image"
	"github.com/metrue/fx/utils"
)

// Docker docker as image builder
type Docker struct {
	*client.Client
}

// CreateClient create a docker instance
func CreateClient() (*Docker, error) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	cli.NegotiateAPIVersion(ctx)
	return &Docker{cli}, nil
}

// Build a directory to be a image
func (d *Docker) Build(workdir string, name string) error {
	tarDir, err := ioutil.TempDir("/tmp", "fx-tar")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tarDir)

	imageID := uuid.New().String()
	tarFilePath := filepath.Join(tarDir, fmt.Sprintf("%s.tar", imageID))

	if err := utils.TarDir(workdir, tarFilePath); err != nil {
		return err
	}

	dockerBuildContext, err := os.Open(tarFilePath)
	if err != nil {
		return err
	}
	defer dockerBuildContext.Close()

	options := dockerTypes.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{imageID, name},
		Labels: map[string]string{
			"belong-to": "fx",
		},
	}

	resp, err := d.ImageBuild(context.Background(), dockerBuildContext, options)
	if err != nil {
		return err
	}

	if os.Getenv("DEBUG") != "" {
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		log.Info(string(body))
	}

	return nil
}

// Push image to hub.docker.com
func (d *Docker) Push(name string) error {
	ctx := context.Background()

	username := os.Getenv("DOCKER_USERNAME")
	password := os.Getenv("DOCKER_PASSWORD")
	if username == "" || password == "" {
		return fmt.Errorf("DOCKER_USERNAME and DOCKER_PASSWORD required for push image to registy")
	}

	authConfig := dockerTypes.AuthConfig{
		Username: username,
		Password: password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return err
	}

	nameWithTag := username + "/" + name
	if err := d.ImageTag(ctx, name, nameWithTag); err != nil {
		return err
	}

	options := dockerTypes.ImagePushOptions{
		RegistryAuth: base64.URLEncoding.EncodeToString(encodedJSON),
	}
	resp, err := d.ImagePush(ctx, nameWithTag, options)
	if err != nil {
		return err
	}
	defer resp.Close()

	if os.Getenv("DEBUG") != "" {
		body, err := ioutil.ReadAll(resp)
		if err != nil {
			return err
		}
		log.Info(string(body))
	}
	return nil
}

var (
	_ image.Builder = &Docker{}
)
