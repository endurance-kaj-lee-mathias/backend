package application

import (
	"context"
	"errors"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/infrastructure/entities"
)

func (s *service) GetOrCreate(ctx context.Context, id domain.UserId, email string, firstName string, lastName string, roles []domain.Role) (domain.User, error) {
	// First try to load the user from the database. If it exists, that record
	// is the source of truth and we ignore the JWT payload.
	usr, err := s.GetByID(ctx, id)
	if err == nil {
		return usr, nil
	}

	// If the error was something other than "not found", bubble it up.
	if !errors.Is(err, infrastructure.NotFound) {
		return domain.User{}, err
	}

	// User does not exist yet: we need at least an email to create the record.
	if email == "" {
		return domain.User{}, errors.New("email required")
	}

	usr = domain.NewUser(id, email, firstName, lastName, roles)
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

func (s *service) DeleteUser(ctx context.Context, id domain.UserId) error {
	return s.repo.Delete(ctx, id.UUID)
}
