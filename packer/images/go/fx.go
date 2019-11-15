package main

import "github.com/gin-gonic/gin"

func fx(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "hello world",
	})
}
