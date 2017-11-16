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
	"cache_dir":        path.Join(os.Getenv("HOME"), ".fx/"),
	"remote_images_url": "https://raw.githubusercontent.com/metrue/fx/master/images.zip",
}
