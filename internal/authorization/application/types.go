package application

import (
	"context"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/authorization/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/authorization/infrastructure"
)

type Service interface {
	CreateRule(ctx context.Context, actorID uuid.UUID, ownerID uuid.UUID, viewerID uuid.UUID, resource string, effect string) (domain.Rule, error)
	DeleteRule(ctx context.Context, actorID uuid.UUID, ruleID uuid.UUID) error
	ListRules(ctx context.Context, ownerID uuid.UUID) ([]domain.Rule, error)
	ListRulesByViewer(ctx context.Context, ownerID uuid.UUID, viewerID uuid.UUID) ([]domain.Rule, error)
	RevokeAll(ctx context.Context, ownerID uuid.UUID, viewerID uuid.UUID) error
	IsAllowed(ctx context.Context, ownerID uuid.UUID, viewerID uuid.UUID, resource string) (bool, error)
	SetResourcePrivacy(ctx context.Context, actorID uuid.UUID, resource string, isPrivate bool) error
	GetResourcePrivacySettings(ctx context.Context, ownerID uuid.UUID) (map[string]bool, error)
	HasSupportRelationship(ctx context.Context, userA uuid.UUID, userB uuid.UUID) (bool, error)
}

type service struct {
	repo infrastructure.Repository
}

func NewService(repo infrastructure.Repository) Service {
	return &service{repo: repo}
}
