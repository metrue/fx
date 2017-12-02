package env

import (
	"fmt"
	"path"

	"github.com/metrue/fx/config"
	api "github.com/metrue/fx/server/docker-api"
	"github.com/metrue/fx/utils"
)

func PullBaseDockerImage(verbose bool) {
	baseImages := []string{
		"metrue/fx-java-base",
		"metrue/fx-julia-base",
		"metrue/fx-python-base",
		"metrue/fx-node-base",
	}

	task := func(image string, verbose bool) {
		fmt.Println("Pulling %s", image)
		api.Pull(image, verbose)
	}

	for _, image := range baseImages {
		go task(image, verbose)
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
func Init(verbose bool) {
	exist, err := utils.IsPathExists(path.Join(config.Client["cache_dir"], "images"))
	if err != nil {
		panic(err)
	}
	if !exist {
		FetchPresetDockerfile()
	}
	PullBaseDockerImage(verbose)
}
