package main

import (
	"path"
	"io/ioutil"
	"path/filepath"
	// "fmt"

	"./common"
	"./utils"
)

func GetGoImageAssets() {
	names := []string{
		"images/go/Dockerfile",
		"images/go/app.go",
		"images/go/fx.go",
	}
	for _, name := range names {
		data, err := common.Asset(name)
		if err != nil {
			panic(err)
		}

		fp := path.Join("/tmp", name)
		dir := filepath.Dir(fp)
		utils.EnsurerDir(dir);

		werr :=ioutil.WriteFile(fp, data, 0644)
		if werr != nil {
			panic(werr)
		}
	}
}

func main() {
	GetGoImageAssets()
}

// images/go/Dockerfile
// images/go/app.go
// images/go/fx.go
// images/java/Dockerfile
// images/java/pom.xml
// images/java/src/main/java/fx/Fx.java
// images/java/src/main/java/fx/app.java
// images/julia/Dockerfile
// images/julia/REQUIRE
// images/julia/app.jl
// images/julia/deps.jl
// images/julia/fx.jl
// images/node/Dockerfile
// images/node/app.js
// images/node/fx.js
// images/php/Dockerfile
// images/php/fx.php
// images/php/index.php
// images/python/Dockerfile
// images/python/app.py
// images/python/fx.py
// images/ruby/Dockerfile
// images/ruby/app.rb
// images/ruby/fx.rb

