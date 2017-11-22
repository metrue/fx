package env

import (
	"fmt"
	"path"

	"github.com/metrue/fx/config"
	"github.com/metrue/fx/utils"
)

// Init creates the server
func Init() {
	exist, err := utils.IsPathExists(path.Join(config.Client["cache_dir"], "images"))
	if err != nil {
		panic(err)
	}
	if !exist {
		fmt.Println("Downloading Resources ...")
		if err := utils.Download("./images.zip", config.Client["remote_images_url"]); err != nil {
			panic(err)
		}
		if err := utils.Unzip("./images.zip", config.Client["cache_dir"]); err != nil {
			panic(err)
		}
	}
}
