package controller

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"github.com/LuisBAndrade/notify/websocket/manager"
	"github.com/LuisBAndrade/notify/websocket/stream"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleWebsocket(c *gin.Context) {
	userID := c.Query("user_id")
	clientName := c.Query("client_name")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user_id"})
		return 
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade websocket: %v", err)
		return 
	}

	manager := websocketmanager.GetManager()
	manager.Connect(conn, userID)
	log.Printf("Websocket connected for user %s (client %s)", userID, clientName)

	ctx := context.Background()

	if err := redisstream.CreateConsumerGroup(ctx, userID); err != nil {
		log.Printf("Failed to create consumer group for user %s: %v", userID, err)
		_ = conn.Close()
		return 
	}

	if pending, err := redisstream.GetPendingNotifications(ctx, userID); err == nil {
		for _, msg := range pending {
			if val, ok := msg.Values["message"].(string); ok {
				payload := map[string]string{
					"message_id": msg.ID,
					"message": val,
				}
				_ = conn.WriteJSON(payload)
			}
		}
	} else {
		log.Printf("Error fetching pending notifications for user %s: %v", userID, err)
	}

	go redisstream.ListenForNotifications(ctx, userID, func(msgID string, message string) error {
		payload := map[string]string{
			"message_id": msgID,
			"message": message,
		}
		return conn.WriteJSON(payload)
	})

	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			log.Printf("WebSocket disconnected for user %s: %v", userID, err)
			break
		}

		text := string(data)
		log.Printf("Received message from %s: %s", userID, text)

		if _, err := redisstream.PublishMessage(ctx, userID, "[Echo] " + text); err != nil {
			log.Printf("Error publishing user message: %v", err)
		}
	}

	manager.Disconnect(conn, userID)
	log.Printf("Connection cleanup done for user %s", userID)
}