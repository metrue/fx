package handlers

import (
	"log"
	"time"
	api "../docker-api"
)

func Down(
	containID string,
	msgChan chan<- string,
	doneChan chan<- bool,
) {
	checkErr := func(err error) bool {
		if err != nil {
			log.Println(err)
			doneChan <- false
			return true
		}
		return false
	}

	err = api.Stop(containID)
	if checkErr(err) {
		return
	}

	fmt.Println("I am closed " + containID)
	msgChan <- containID + " Stopped"
	doneChan <- true
}
