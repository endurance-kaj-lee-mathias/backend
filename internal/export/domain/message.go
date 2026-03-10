package domain

import (
	"time"
)

type MessageData struct {
	ID                        string    `json:"id"`
	ConversationID            string    `json:"conversationId"`
	SenderID                  string    `json:"senderId"`
	Content                   string    `json:"content"`
	CreatedAt                 time.Time `json:"createdAt"`
	OtherParticipantUsername  string    `json:"otherParticipantUsername"`
	OtherParticipantFirstName string    `json:"otherParticipantFirstName"`
	OtherParticipantLastName  string    `json:"otherParticipantLastName"`
}
