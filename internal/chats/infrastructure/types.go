package infrastructure

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/chats/infrastructure/entities"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
)

type Repository interface {
	FindConversation(ctx context.Context, userA, userB uuid.UUID) (entities.ConversationEntity, error)
	FindConversations(ctx context.Context, userID uuid.UUID) ([]entities.ConversationWithParticipantsEntity, error)
	CreateConversation(ctx context.Context, ent entities.ConversationEntity) error
	SaveParticipantKey(ctx context.Context, ent entities.ParticipantKeyEntity) error
	UpdateParticipantKey(ctx context.Context, conversationID, userID uuid.UUID, encryptedKey []byte) error
	GetParticipantKey(ctx context.Context, conversationID, userID uuid.UUID) (entities.ParticipantKeyEntity, error)
	GetUserEncryptedKey(ctx context.Context, userID uuid.UUID) ([]byte, error)
	CreateMessage(ctx context.Context, ent entities.MessageEntity) error
	GetMessages(ctx context.Context, conversationID uuid.UUID, limit, offset int) ([]entities.MessageEntity, error)
	CheckSupportRelationship(ctx context.Context, userA, userB uuid.UUID) (bool, error)
}

type repository struct {
	db  *sql.DB
	enc encryption.Service
}

func NewRepository(db *sql.DB, enc encryption.Service) Repository {
	return &repository{db: db, enc: enc}
}
