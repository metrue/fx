package image

import (
	"github.com/metrue/common"
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

var assetsMap = map[string][string] {
	"go": {
		"assets/images/go/Dockerfile",
		"assets/images/go/app.go",
		"assets/images/go/fx.go",
	},
	"java": {
		"assets/dockerfiles/fx/java/Dockerfile",
		"assets/dockerfiles/fx/java/pom.xml",
		"assets/dockerfiles/fx/java/src/main/java/fx/Fx.java"
		"assets/dockerfiles/fx/java/src/main/java/fx/app.java",
	},
	"julia": {
		"assets/dockerfiles/fx/julia/Dockerfile",
		"assets/dockerfiles/fx/julia/REQUIRE",
		"assets/dockerfiles/fx/julia/app.jl",
		"assets/dockerfiles/fx/julia/deps.jl"
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
	}
}

func Get(dir string, lang []byte, body []byte) (err error){
	names := assetsMap[lang]
	for _, name := range names {
		data, err := common.Asset(name);

		targetPath := path.Join(dir,name)
		err := ioutil.WriteFile(targetPath, data, 0644)
	}
}
