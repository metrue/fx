package env

import (
	"../../utils"
	Config "../../config"
	"fmt"
	"path"
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
