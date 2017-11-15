package config

import (
	"path"
	"os"
	"flag"
)

var CacheDir = path.Join(os.Getenv("HOME"), ".fx/")
var RemoteImagesUrl = "https://raw.githubusercontent.com/metrue/fx/master/images.zip"
var ServerAddr = flag.String("addr", "localhost:8080", "http service address")
