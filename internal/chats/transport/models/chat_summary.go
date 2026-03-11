package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/chats/domain"
)

type ChatSummaryModel struct {
	ConversationID      uuid.UUID  `json:"conversationId"`
	OtherUserID         uuid.UUID  `json:"otherUserId"`
	FirstName           string     `json:"firstName"`
	LastName            string     `json:"lastName"`
	ImageUrl            string     `json:"imageUrl"`
	LatestMessage       *string    `json:"latestMessage"`
	LatestMessageSentBy *uuid.UUID `json:"latestMessageSentBy"`
	LatestMessageAt     *time.Time `json:"latestMessageAt"`
}

func ToChatSummaryModel(s domain.ConversationSummary) ChatSummaryModel {
	return ChatSummaryModel{
		ConversationID:      s.ConversationID.UUID,
		OtherUserID:         s.OtherUserID,
		FirstName:           s.FirstName,
		LastName:            s.LastName,
		ImageUrl:            s.Image,
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
