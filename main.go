package main

import (
	"io/ioutil"
	"path"
	"path/filepath"
	// "fmt"

	"github.com/metrue/fx/common"
	"github.com/metrue/fx/utils"
)

func GetGoImageAssets() {
	names := []string{
		"assets/dockerfiles/fx/go/Dockerfile",
		"assets/dockerfiles/fx/go/app.go",
		"assets/dockerfiles/fx/go/fx.go",
	}
	for _, name := range names {
		data, err := common.Asset(name)
		if err != nil {
			panic(err)
		}

		fp := path.Join("/tmp", name)
		dir := filepath.Dir(fp)
		utils.EnsurerDir(dir)

		werr := ioutil.WriteFile(fp, data, 0644)
		if werr != nil {
			panic(werr)
		}
	}
}

func main() {
	GetGoImageAssets()
}
