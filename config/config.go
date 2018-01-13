package config

import (
	"os"
	"path"
)

// Server contains the server configuration information
var Server = map[string]string{
	"host": "localhost",
	"port": "30080",
}

// Client contains the local and remote paths to fetch cached images
var Client = map[string]string{
	"cache_dir":         path.Join(os.Getenv("HOME"), ".fx/"),
	"remote_images_url": "https://raw.githubusercontent.com/metrue/fx/master/images.zip",
}

var GrpcEndpoint = ":5000"
