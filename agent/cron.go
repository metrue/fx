package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

func request(endpoint string, params map[string]string) error {
	body, err := json.Marshal(params)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func enableCron(ctx *gin.Context) {
	type Payload struct {
		Cron       string `json:"cron" form:"cron"`
		Service    string `json:"service" form:"service"`
		Parameters string `json:"parameters" form:"parameters"`
	}

	var payload Payload
	if err := ctx.ShouldBind(&payload); err != nil {
		ctx.JSON(400, gin.H{
			"error": fmt.Sprintf("invalid request payload: %s", err),
		})
		return
	}

	var params map[string]string
	groups := strings.Split(payload.Parameters, ",")
	for _, g := range groups {
		param := strings.Split(g, "=")
		if len(param) == 2 {
			params[param[0]] = param[1]
		} else {
			fmt.Println(g, param)
		}
	}

	c := cron.New()
	fmt.Println("cron is ", payload.Cron)
	err := c.AddFunc(payload.Cron, func() {
		fmt.Println("call")
		if err := request(payload.Service, params); err != nil {
			fmt.Println("called failed", err)
		} else {
			fmt.Println("called ok")
		}
	})
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": fmt.Sprintf("cron enable failed: %s", err),
		})
		return
	}
	c.Start()
	ctx.JSON(201, gin.H{
		"message": "cron enabled",
	})
}

func main() {
	r := gin.Default()

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"msg": "pong"})
	})
	r.POST("/ping", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"msg": "pong"})
	})
	r.POST("/cron", enableCron)

	if err := r.Run(":3000"); err != nil {
		panic(err)
	}
}
