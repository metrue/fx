package config

import (
	"fmt"
	"os"
	"path"
)

func Config() {
	address := os.Args[2]
	home := os.Getenv("HOME")
	file := path.Join(home, ".fx")
	f, err := os.Create(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	f.WriteString(address)
}
