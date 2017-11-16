package env

import (
	"fmt"
	"path"

	Config "github.com/metrue/fx/config"
	"github.com/metrue/fx/utils"
)

func Init() {
	exist, err := utils.IsPathExists(path.Join(Config.Client["cache_dir"], "images"))
	if err != nil {
		panic(err)
	}
	if !exist {
		fmt.Println("Downloading Resources ...")
		if err := utils.Download("./images.zip", Config.Client["remote_images_url"]); err != nil {
			panic(err)
		}
		if err := utils.Unzip("./images.zip", Config.Client["cache_dir"]); err != nil {
			panic(err)
		}
	}
}
