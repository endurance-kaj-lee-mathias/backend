package application

import (
	"context"
	"errors"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/infrastructure/entities"
)

func (s *service) GetOrCreate(ctx context.Context, id domain.UserId, email string, firstName string, lastName string, roles []domain.Role) (domain.User, error) {
	usr, err := s.GetByID(ctx, id)
	if err == nil {
		return usr, nil
	}

	if !errors.Is(err, infrastructure.NotFound) {
		return domain.User{}, err
	}

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

func (s *service) UpdatePhoneNumber(ctx context.Context, id domain.UserId, phoneNumber *string) error {
	return s.repo.UpdatePhoneNumber(ctx, id.UUID, phoneNumber)
}

func (s *service) UpsertAddress(ctx context.Context, userID domain.UserId, street string, houseNumber string, postalCode string, city string, country string) (domain.Address, error) {
	addrID, err := domain.NewAddressId()
	if err != nil {
		return domain.Address{}, err
	}

	addr := domain.NewAddress(addrID, userID, street, houseNumber, postalCode, city, country)
	ent := entities.AddressToEntity(addr)

	if err := s.repo.InsertAddress(ctx, ent); err != nil {
		return domain.Address{}, err
	}

	return s.GetAddress(ctx, userID)
}

func (s *service) GetAddress(ctx context.Context, userID domain.UserId) (domain.Address, error) {
	ent, err := s.repo.FindAddressByUserID(ctx, userID.UUID)
	if err != nil {
		return domain.Address{}, err
	}

	return entities.AddressFromEntity(ent)
}
