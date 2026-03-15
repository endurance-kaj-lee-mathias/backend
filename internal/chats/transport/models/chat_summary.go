package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/chats/domain"
)

type ChatSummaryModel struct {
	ID                  uuid.UUID  `json:"id"`
	OtherUserID         uuid.UUID  `json:"otherUserId"`
	Username            string     `json:"username"`
	FirstName           string     `json:"firstName"`
	LastName            string     `json:"lastName"`
	Image               string     `json:"image"`
	LatestMessage       *string    `json:"latestMessage"`
	LatestMessageSentBy *uuid.UUID `json:"latestMessageSentBy"`
	LatestMessageAt     *time.Time `json:"latestMessageAt"`
}

func ToChatSummaryModel(s domain.ConversationSummary) ChatSummaryModel {
	return ChatSummaryModel{
		ID:                  s.ConversationID.UUID,
		OtherUserID:         s.OtherUserID,
		Username:            s.Username,
		FirstName:           s.FirstName,
		LastName:            s.LastName,
		Image:               s.Image,
		LatestMessage:       s.LatestMessage,
		LatestMessageSentBy: s.LatestMessageSenderID,
		LatestMessageAt:     s.LatestMessageAt,
	}
}

func ToChatSummaryModels(summaries []domain.ConversationSummary) []ChatSummaryModel {
	out := make([]ChatSummaryModel, 0, len(summaries))
	for _, s := range summaries {
		out = append(out, ToChatSummaryModel(s))
	}
	return out
}
