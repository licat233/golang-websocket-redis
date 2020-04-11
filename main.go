package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang-websocket-redis/redis"
	"golang-websocket-redis/server"
)

func main() {
	go redis.SetMsgLoop()
	R := gin.Default()
	R.GET("/ws", func(c *gin.Context) {
		server.WebsocketHandler(c.Writer, c.Request)
	})
	if err := R.Run(":8888"); err != nil {
		fmt.Println("Gin run failed")
		return
	}
}
