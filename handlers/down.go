package handlers

import (
	"fmt"
	"log"

	api "github.com/metrue/fx/docker-api"
)

// Down stops the processes designated by a function
func Down(containID, image string, msgChan chan<- string, doneChan chan<- bool) {
	checkErr := func(err error) bool {
		if err != nil {
			log.Println(err)
			doneChan <- false
			return true
		}
		return false
	}

	err := api.Remove(containID)
	if checkErr(err) {
		return
	}

	fmt.Println("I am closed " + containID)
	msgChan <- fmt.Sprintf("Container[%s] Removed", containID)
	if err := api.ImageRemove(image); err != nil {
		log.Printf("cleanup docker image[%s] error: %s\n", image, err.Error())
		msgChan <- fmt.Sprintf("Image[%s] Removing Failed", image)
	} else {
		msgChan <- fmt.Sprintf("Image[%s] Removed", image)
	}

	doneChan <- true
}
