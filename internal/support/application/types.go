package application

import (
	"context"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/infrastructure"
)

type Service interface {
	GetAll(ctx context.Context, id domain.VeteranId) ([]domain.Member, error)
	GetAllByMember(ctx context.Context, id domain.MemberId) ([]domain.Member, error)
	DeleteSupporter(ctx context.Context, veteranID domain.VeteranId, supportID domain.MemberId) error

	SendInvite(ctx context.Context, senderID domain.MemberId, receiverID domain.MemberId) (domain.Invite, error)
	AcceptInvite(ctx context.Context, callerID domain.MemberId, inviteID domain.InviteId) (domain.Invite, error)
	DeclineInvite(ctx context.Context, callerID domain.MemberId, inviteID domain.InviteId) error
	ListInvites(ctx context.Context, callerID domain.MemberId) (incoming []domain.Invite, outgoing []domain.Invite, err error)
}

type service struct {
	repo         infrastructure.Repository
	inviteRepo   infrastructure.InviteRepository
	userRoleRead infrastructure.UserRoleReader
	enc          encryption.Service
}

func NewService(repo infrastructure.Repository, inviteRepo infrastructure.InviteRepository, userRoleRead infrastructure.UserRoleReader, enc encryption.Service) Service {
	return &service{repo: repo, inviteRepo: inviteRepo, userRoleRead: userRoleRead, enc: enc}
}
