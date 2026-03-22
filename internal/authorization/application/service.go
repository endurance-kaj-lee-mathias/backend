package application

import (
	"context"
	"errors"
	"log/slog"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/authorization/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/authorization/infrastructure"
)

func (s *service) CreateRule(ctx context.Context, actorID uuid.UUID, ownerID uuid.UUID, viewerID uuid.UUID, resource string, effect string) (domain.Rule, error) {
	if actorID != ownerID {
		return domain.Rule{}, domain.NotOwner
	}

	if ownerID == viewerID {
		return domain.Rule{}, domain.SelfRule
	}

	if !domain.ValidResource(resource) {
		return domain.Rule{}, domain.InvalidResource
	}

	if !domain.ValidEffect(effect) {
		return domain.Rule{}, domain.InvalidEffect
	}

	rule, err := domain.NewRule(ownerID, viewerID, domain.ResourceType(resource), domain.PolicyEffect(effect))
	if err != nil {
		return domain.Rule{}, err
	}

	if err := s.repo.Create(ctx, rule); err != nil {
		return domain.Rule{}, err
	}

	slog.Info("authorization rule created",
		"owner", ownerID.String(),
		"viewer", viewerID.String(),
		"resource", resource,
		"effect", effect,
	)

	return rule, nil
}

func (s *service) DeleteRule(ctx context.Context, actorID uuid.UUID, ruleID uuid.UUID) error {
	rule, err := s.repo.FindByID(ctx, ruleID)
	if err != nil {
		if errors.Is(err, infrastructure.RuleNotFound) {
			return domain.RuleNotFound
		}
		return err
	}

	if rule.OwnerID != actorID {
		return domain.NotOwner
	}

	if err := s.repo.Delete(ctx, ruleID); err != nil {
		return err
	}

	slog.Info("authorization rule deleted",
		"owner", rule.OwnerID.String(),
		"viewer", rule.ViewerID.String(),
		"resource", string(rule.Resource),
		"ruleId", ruleID.String(),
	)

	return nil
}

func (s *service) ListRules(ctx context.Context, ownerID uuid.UUID) ([]domain.Rule, error) {
	rules, err := s.repo.FindByOwner(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	if rules == nil {
		rules = []domain.Rule{}
	}

	return rules, nil
}

func (s *service) ListRulesByViewer(ctx context.Context, ownerID uuid.UUID, viewerID uuid.UUID) ([]domain.Rule, error) {
	rules, err := s.repo.FindByOwnerAndViewer(ctx, ownerID, viewerID)
	if err != nil {
		return nil, err
	}

	if rules == nil {
		rules = []domain.Rule{}
	}

	return rules, nil
}

func (s *service) RevokeAll(ctx context.Context, ownerID uuid.UUID, viewerID uuid.UUID) error {
	if err := s.repo.DeleteByOwnerAndViewer(ctx, ownerID, viewerID); err != nil {
		return err
	}

	slog.Info("all authorization rules revoked",
		"owner", ownerID.String(),
		"viewer", viewerID.String(),
	)

	return nil
}

func (s *service) IsAllowed(ctx context.Context, ownerID uuid.UUID, viewerID uuid.UUID, resource string) (bool, error) {
	if ownerID == viewerID {
		return true, nil
	}

	rule, err := s.repo.FindRule(ctx, ownerID, viewerID, resource)
	if err != nil {
		return false, err
	}

	if rule != nil {
		return rule.Effect == domain.EffectAllow, nil
	}

	resourcePrivate, err := s.repo.GetResourcePrivacy(ctx, ownerID, resource)
	if err != nil {
		return false, err
	}

	if resourcePrivate != nil {
		return !*resourcePrivate, nil
	}

	isPrivate, err := s.repo.GetPrivacy(ctx, ownerID)
	if err != nil {
		return false, err
	}

	return !isPrivate, nil
}

func (s *service) SetResourcePrivacy(ctx context.Context, actorID uuid.UUID, resource string, isPrivate bool) error {
	if !domain.ValidResource(resource) {
		return domain.InvalidResource
	}

	if err := s.repo.SetResourcePrivacy(ctx, actorID, resource, isPrivate); err != nil {
		return err
	}

	slog.Info("resource privacy updated",
		"owner", actorID.String(),
		"resource", resource,
		"isPrivate", isPrivate,
	)

	return nil
}

func (s *service) GetResourcePrivacySettings(ctx context.Context, ownerID uuid.UUID) (map[string]bool, error) {
	return s.repo.ListResourcePrivacy(ctx, ownerID)
}

func (s *service) HasSupportRelationship(ctx context.Context, userA uuid.UUID, userB uuid.UUID) (bool, error) {
	return s.repo.HasSupportRelationship(ctx, userA, userB)
}
