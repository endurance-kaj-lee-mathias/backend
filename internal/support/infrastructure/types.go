package infrastructure

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/infrastructure/entities"
)

type UserRoleReader interface {
	GetRole(ctx context.Context, userID uuid.UUID) (string, error)
	FindIDByUsername(ctx context.Context, username string) (uuid.UUID, error)
}

type Repository interface {
	Create(ctx context.Context, veteranID, memberId uuid.UUID) (entities.MemberEntity, error)
	ReadAll(ctx context.Context, id uuid.UUID) ([]entities.MemberEntity, error)
	ReadAllByMember(ctx context.Context, id uuid.UUID) ([]entities.MemberEntity, error)
	Delete(ctx context.Context, veteranID, supportID uuid.UUID) error
}

type InviteRepository interface {
	CreateInvite(ctx context.Context, inv domain.Invite) error
	FindInviteByID(ctx context.Context, id uuid.UUID) (entities.InviteEntity, error)
	FindPendingBySenderReceiver(ctx context.Context, senderID, receiverID uuid.UUID) (entities.InviteEntity, bool, error)
	FindAcceptedBySenderReceiver(ctx context.Context, senderID, receiverID uuid.UUID) (bool, error)
	UpdateInviteStatus(ctx context.Context, id uuid.UUID, status domain.InviteStatus) error
	DeleteInvite(ctx context.Context, id uuid.UUID) error
	ListPendingForUser(ctx context.Context, userID uuid.UUID) ([]entities.InviteEntity, error)
}

type userRoleReader struct {
	db  *sql.DB
	enc encryption.Service
}

func NewUserRoleReader(db *sql.DB, enc encryption.Service) UserRoleReader {
	return &userRoleReader{db: db, enc: enc}
}

type repository struct {
	db  *sql.DB
	enc encryption.Service
}

func NewRepository(db *sql.DB, enc encryption.Service) Repository {
	return &repository{db: db, enc: enc}
}

type inviteRepository struct {
	db  *sql.DB
	enc encryption.Service
}

func NewInviteRepository(db *sql.DB, enc encryption.Service) InviteRepository {
	return &inviteRepository{db: db, enc: enc}
}
