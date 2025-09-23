package models

type NotificationRequestData struct {
	UserID string `json:"user_id" binding:"required"`
	Payload map[string]interface{} `json:"payload"`
}

type NotificationResponse struct {
	Success bool `json:"success"`
	MessageID string `json:"message_id"`
}

type AcknowledgeRequest struct {
	UserID string `json:"user_id" binding:"required"`
	MessageIDs []string `json:"message_ids"`
}

type AcknowledgeResponse struct {
	Success bool `json:"success"`
}