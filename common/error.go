package common

import (
	"fmt"
	"os"
)

func HandleError(err error) {
	panic(err)
}

func HandleEnvError(err error) {
	fmt.Println(err.Error())
	os.Exit(0)
}
