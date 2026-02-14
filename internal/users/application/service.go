package application

import (
	"context"
	"errors"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/infrastructure/entities"
)

func (s *service) AddUser(ctx context.Context, email string, roles []domain.Role) (domain.User, error) {
	if email == "" {
		return domain.User{}, errors.New("email required")
	}

	usr := domain.NewUser(email, roles)
	ent, err := entities.ToEntity(usr)

	if err != nil {
		return domain.User{}, err
	}

	if err := s.repo.Create(ctx, ent); err != nil {
		return domain.User{}, err
	}

	return usr, nil
}

func (s *service) GetByID(ctx context.Context, id domain.UserId) (domain.User, error) {
	ent, err := s.repo.FindByID(ctx, id.UUID)

	if err != nil {
		return domain.User{}, err
	}

	usr, err := entities.FromEntity(ent)

	if err != nil {
		return domain.User{}, err
	}

	return usr, nil
}

func (s *service) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	ent, err := s.repo.FindByEmail(ctx, email)

	if err != nil {
		return domain.User{}, err
	}

	usr, err := entities.FromEntity(ent)

	if err != nil {
		return domain.User{}, err
	}

	return usr, nil
}
