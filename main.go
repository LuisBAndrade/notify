package main

import (

	"github.com/LuisBAndrade/notify/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.StaticFile("/", "./client.html")

	router.GET("/api/ws", controller.HandleWebsocket)
	router.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	router.POST("/api/notification/send", controller.SendNotification)
	router.POST("/api/notification/acknowledge", controller.AcknowledgeNotification)

	router.Run(":8080")
}