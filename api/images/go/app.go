package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", fx)
	r.POST("/", fx)
	if err := r.Run(":3000"); err != nil {
		panic(err)
	}
}
