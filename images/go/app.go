package main

import (
	"encoding/json"
	"net/http"
)

func handler(res http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)

	var input Input
	err := decoder.Decode(&input)
	if err != nil {
		panic(err)
	}

	output := Fx(&input)

	encoder := json.NewEncoder(res)
	err = encoder.Encode(output)
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":3000", nil)
}
