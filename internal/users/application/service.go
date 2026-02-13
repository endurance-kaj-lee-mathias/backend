package application

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/infrastructure/entities"
)

func (s *service) CreateUser(ctx context.Context, email string, roles []domain.Role) (*domain.User, error) {
	if email == "" {
		return nil, errors.New("email required")
	}

	user := domain.NewUser(email, roles)

	ent, err := entities.ToEntity(user)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Save(ctx, ent); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *service) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	ent, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	user, err := entities.FromEntity(ent)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *service) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	ent, err := s.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	user, err := entities.FromEntity(ent)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *service) AddSupportMember(ctx context.Context, veteranID, supportID uuid.UUID) error {
	return s.repo.AddSupportMember(ctx, veteranID, supportID)
}

func (s *service) ListSupportMembers(ctx context.Context, veteranID uuid.UUID) ([]domain.User, error) {
	ents, err := s.repo.ListSupportMembers(ctx, veteranID)
	if err != nil {
		return nil, err
	}

	out := make([]domain.User, 0, len(ents))
	for _, ent := range ents {
		user, err := entities.FromEntity(ent)
		if err != nil {
			return nil, err
		}
		out = append(out, user)
	}

	return out, nil
}
