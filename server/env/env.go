package env

import (
	"fmt"
	"path"

	"github.com/metrue/fx/config"
	api "github.com/metrue/fx/server/docker-api"
	"github.com/metrue/fx/utils"
)

func PullBaseDockerImage() {
	baseImages := []string{
		"metrue/fx-java-base",
		"metrue/fx-julia-base",
		"metrue/fx-python-base",
		"metrue/fx-node-base",
	}
	verbose := true
	for _, image := range baseImages {
		api.Pull(image, verbose)
	}
}

func FetchPresetDockerfile() {
	fmt.Println("Downloading Resources ...")
	if err := utils.Download("./images.zip", config.Client["remote_images_url"]); err != nil {
		panic(err)
	}
	if err := utils.Unzip("./images.zip", config.Client["cache_dir"]); err != nil {
		panic(err)
	}
}

// Init creates the server
func Init() {
	exist, err := utils.IsPathExists(path.Join(config.Client["cache_dir"], "images"))
	if err != nil {
		panic(err)
	}
	if !exist {
		FetchPresetDockerfile()
	}
	PullBaseDockerImage()
}
