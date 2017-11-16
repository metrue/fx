package config

import (
	"os"
	"path"
)

const Server = map[string]string{
	"host": "localhost",
	"port": "8080",
}

const Client = map[string]string{
	"CacheDir":        path.Join(os.Getenv("HOME"), ".fx/"),
	"RemoteImagesUrl": "https://raw.githubusercontent.com/metrue/fx/master/images.zip",
}
