package application

import (
	"context"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/infrastructure/entities"
)

func (s *service) AddMember(ctx context.Context, veteranID domain.VeteranId, memberId domain.MemberId) (domain.Member, error) {
	veteranRoles, err := s.userRoleRead.GetRoles(ctx, veteranID.UUID)
	if err != nil {
		return domain.Member{}, err
	}

	supporterRoles, err := s.userRoleRead.GetRoles(ctx, memberId.UUID)
	if err != nil {
		return domain.Member{}, err
	}

	if err := domain.ValidateSupportRelationship(veteranRoles, supporterRoles, veteranID.UUID.String(), memberId.UUID.String()); err != nil {
		return domain.Member{}, err
	}

	ent, err := s.repo.Create(ctx, veteranID.UUID, memberId.UUID)
	if err != nil {
		return domain.Member{}, err
	}

	return entities.FromEntity(ent)
}

func (s *service) GetAll(ctx context.Context, id domain.VeteranId) ([]domain.Member, error) {
	ents, err := s.repo.ReadAll(ctx, id.UUID)

	if err != nil {
		return nil, err
	}

	return entities.FromEntities(ents), nil
}

func (s *service) GetAllByMember(ctx context.Context, id domain.MemberId) ([]domain.Member, error) {
	ents, err := s.repo.ReadAllByMember(ctx, id.UUID)

	if err != nil {
		return nil, err
	}

	return entities.FromEntities(ents), nil
}

func (s *service) DeleteSupporter(ctx context.Context, veteranID domain.VeteranId, supportID domain.MemberId) error {
	return s.repo.Delete(ctx, veteranID.UUID, supportID.UUID)
}
