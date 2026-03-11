package application

import (
	"context"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/chats/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/chats/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
)

type Service interface {
	GetOrCreateConversation(ctx context.Context, callerID, participantID uuid.UUID) (domain.Conversation, error)
	SendMessage(ctx context.Context, conversationID uuid.UUID, senderID uuid.UUID, content string) (domain.Message, error)
	GetMessages(ctx context.Context, conversationID uuid.UUID, callerID uuid.UUID, limit, offset int) ([]domain.Message, error)
	GetAllChats(ctx context.Context, userID uuid.UUID) ([]domain.ConversationSummary, error)
}

type service struct {
	repo infrastructure.Repository
	enc  encryption.Service
}

func NewService(repo infrastructure.Repository, enc encryption.Service) Service {
	return &service{repo: repo, enc: enc}
}
