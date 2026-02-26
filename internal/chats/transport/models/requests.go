package models

import "github.com/gofrs/uuid"

type StartConversationRequest struct {
	ParticipantID uuid.UUID `json:"participantId"`
}

type SendMessageRequest struct {
	Content string `json:"content"`
}
