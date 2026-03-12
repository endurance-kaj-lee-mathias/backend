package application

import (
	"context"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/infrastructure"
)

type AuthzRevoker interface {
	RevokeAll(ctx context.Context, ownerID uuid.UUID, viewerID uuid.UUID) error
}

type Notifier interface {
	NotifyInvite(ctx context.Context, deviceToken string) error
	NotifyInviteAccepted(ctx context.Context, deviceToken string) error
}

type Service interface {
	GetAll(ctx context.Context, id domain.VeteranId) ([]domain.Member, error)
	GetAllByMember(ctx context.Context, id domain.MemberId) ([]domain.Member, error)
	DeleteSupporter(ctx context.Context, callerID domain.MemberId, otherID domain.MemberId) error

	SendInvite(ctx context.Context, senderID domain.MemberId, username string, note *string) (domain.Invite, error)
	AcceptInvite(ctx context.Context, callerID domain.MemberId, inviteID domain.InviteId) error
	DeclineInvite(ctx context.Context, callerID domain.MemberId, inviteID domain.InviteId) error
	ListInvites(ctx context.Context, callerID domain.MemberId) (incoming []domain.Invite, outgoing []domain.Invite, err error)
}

type service struct {
	repo         infrastructure.Repository
	inviteRepo   infrastructure.InviteRepository
	userRoleRead infrastructure.UserRoleReader
	enc          encryption.Service
	authz        AuthzRevoker
	notifier     Notifier
}

func NewService(repo infrastructure.Repository, inviteRepo infrastructure.InviteRepository, userRoleRead infrastructure.UserRoleReader, enc encryption.Service, authz AuthzRevoker, notifier Notifier) Service {
	return &service{repo: repo, inviteRepo: inviteRepo, userRoleRead: userRoleRead, enc: enc, authz: authz, notifier: notifier}
}
