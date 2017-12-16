package image

import (
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/metrue/fx/common"
	"github.com/metrue/fx/utils"
)

var funcNames = map[string]string{
	"go":     "/fx.go",
	"node":   "/fx.js",
	"ruby":   "/fx.rb",
	"python": "/fx.py",
	"php":    "/fx.php",
	"julia":  "/fx.jl",
	"java":   "/src/main/java/fx/Fx.java",
}

var assetsMap = map[string][]string{
	"go": {
		"assets/dockerfiles/fx/go/Dockerfile",
		"assets/dockerfiles/fx/go/app.go",
		"assets/dockerfiles/fx/go/fx.go",
	},
	"java": {
		"assets/dockerfiles/fx/java/Dockerfile",
		"assets/dockerfiles/fx/java/pom.xml",
		"assets/dockerfiles/fx/java/src/main/java/fx/Fx.java",
		"assets/dockerfiles/fx/java/src/main/java/fx/app.java",
	},
	"julia": {
		"assets/dockerfiles/fx/julia/Dockerfile",
		"assets/dockerfiles/fx/julia/REQUIRE",
		"assets/dockerfiles/fx/julia/app.jl",
		"assets/dockerfiles/fx/julia/deps.jl",
		"assets/dockerfiles/fx/julia/fx.jl",
	},
	"node": {
		"assets/dockerfiles/fx/node/Dockerfile",
		"assets/dockerfiles/fx/node/app.js",
		"assets/dockerfiles/fx/node/fx.js",
	},
	"php": {
		"assets/dockerfiles/fx/php/Dockerfile",
		"assets/dockerfiles/fx/php/fx.php",
		"assets/dockerfiles/fx/php/index.php",
	},
	"python": {
		"assets/dockerfiles/fx/python/Dockerfile",
		"assets/dockerfiles/fx/python/app.py",
		"assets/dockerfiles/fx/python/fx.py",
	},
	"ruby": {
		"assets/dockerfiles/fx/ruby/Dockerfile",
		"assets/dockerfiles/fx/ruby/app.rb",
		"assets/dockerfiles/fx/ruby/fx.rb",
	},
}

func removePrefix(lang string, filename string) (name string) {
	prefix := "assets/dockerfiles/fx" + "/" + lang + "/"
	return strings.Split(filename, prefix)[1]
}

func isFxFuncSource(lang string, name string) (ret bool) {
	basename := filepath.Base(name)
	nameWithoutExt := strings.TrimSuffix(basename, filepath.Ext(basename))
	return nameWithoutExt == "fx" || nameWithoutExt == "Fx" // Fx is for Java
}

//Get Prepare a container base image and insert the function body
func Get(dir string, lang string, body []byte) (err error) {
	names := assetsMap[lang]
	err = nil
	for _, name := range names {
		data, assetErr := common.Asset(name)
		if assetErr != nil {
			err = assetErr
		}

		name = removePrefix(lang, name)
		targetPath := path.Join(dir, name)

		dir := filepath.Dir(targetPath)
		utils.EnsurerDir(dir)

		if isFxFuncSource(lang, targetPath) {
			data = body
		}

		writeErr := ioutil.WriteFile(targetPath, data, 0644)
		if writeErr != nil {
			err = writeErr
		}
	}
	return err
}
