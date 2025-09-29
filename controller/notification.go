package controller

import (
	"context"
	"net/http"

	"github.com/LuisBAndrade/notify/models"
	"github.com/LuisBAndrade/notify/websocket/stream"
	"github.com/gin-gonic/gin"
)

func SendNotification(c *gin.Context) {
	var req models.NotificationRequestData
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return 
	}

	ctx := context.Background()
	msgID, err := redisstream.PublishEvent(ctx, req.UserID, "NOTIFICATION", req.Payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to publish"})
		return 
	}

	c.JSON(http.StatusOK, models.NotificationResponse{
		Success:  true,
		MessageID: msgID,
	})
}

func AcknowledgeNotification(c *gin.Context) {
	var req models.AcknowledgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := context.Background()
	if err := redisstream.AcknowledgeNotifications(ctx, req.UserID, req.MessageIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to acknowledge"})
		return
	}

	_, _ = redisstream.PublishEvent(ctx, req.UserID, "ACK", map[string]interface{}{
		"acknowledged_ids": req.MessageIDs,
	})

	c.JSON(http.StatusOK, models.AcknowledgeResponse{Success: true})
}