package env

import (
	"fmt"
	"path"

	Config "github.com/metrue/fx/config"
	"github.com/metrue/fx/utils"
)

func Init() {
	exist, err := utils.IsPathExists(path.Join(Config.CacheDir, "images"))
	if err != nil {
		panic(err)
	}
	if !exist {
		fmt.Println("Downloading Resources ...")
		if err := utils.Download("./images.zip", Config.RemoteImagesUrl); err != nil {
			panic(err)
		}
		if err := utils.Unzip("./Images.zip", Config.CacheDir); err != nil {
			panic(err)
		}
	}
}
