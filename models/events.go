package models

type EventMessage struct {
	Event string `json:"event"`
	MessageID string `json:"message_id"`
	Payload interface{} `json:"payload"`
}