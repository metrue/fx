package docker

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/apex/log"
	dockerTypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/metrue/fx/builders"
	"github.com/metrue/fx/utils"
)

// Docker docker as image builder
type Docker struct {
	*client.Client
}

// CreateDocker create a docker instance
func CreateDocker() (*Docker, error) {
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

var (
	_ builders.Builder = &Docker{}
)
