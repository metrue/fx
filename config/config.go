package config

import (
	"os"
	"path"
)

var Server = map[string]string{
	"host": "localhost",
	"port": "8080",
}

var Client = map[string]string{
	"CacheDir":        path.Join(os.Getenv("HOME"), ".fx/"),
	"RemoteImagesUrl": "https://raw.githubusercontent.com/metrue/fx/master/images.zip",
}
