package application

import (
	"context"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/chats/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/chats/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/ws/application"
)

type Notifier interface {
	NotifyNewMessage(ctx context.Context, deviceToken string) error
}

type Service interface {
	GetOrCreateConversation(ctx context.Context, callerID, participantID uuid.UUID) (domain.Conversation, error)
	SendMessage(ctx context.Context, conversationID uuid.UUID, senderID uuid.UUID, content string) (domain.Message, error)
	GetMessages(ctx context.Context, conversationID uuid.UUID, callerID uuid.UUID, limit, offset int) ([]domain.Message, error)
	GetAllChats(ctx context.Context, userID uuid.UUID) ([]domain.ConversationSummary, error)
	ListConversationIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}

type service struct {
	repo        infrastructure.Repository
	enc         encryption.Service
	notifier    Notifier
	broadcaster application.Broadcaster
}

func NewService(repo infrastructure.Repository, enc encryption.Service, notifier Notifier, broadcaster application.Broadcaster) Service {
	return &service{repo: repo, enc: enc, notifier: notifier, broadcaster: broadcaster}
}
