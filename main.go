package main

import (

	"github.com/LuisBAndrade/notify/controller"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/api/ws", controller.HandleWebsocket)

	router.POST("/api/notification/send", controller.SendNotification)
	router.POST("/api/notification/acknowledge", controller.AcknowledgeNotification)

	router.Run(":8080")
}