package handlers

import (
	"fmt"
	"log"

	api "github.com/metrue/fx/server/docker-api"
)

// Down stops the processes designated by a function
func Down(containID string, msgChan chan<- string, doneChan chan<- bool) {
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
	msgChan <- containID + " Stopped"
	doneChan <- true
}
