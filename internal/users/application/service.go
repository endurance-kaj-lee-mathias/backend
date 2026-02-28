package application

import (
	"context"
	"errors"
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/infrastructure"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/infrastructure/entities"
)

func (s *service) GetOrCreate(ctx context.Context, id domain.UserId, email string, username string, firstName string, lastName string, roles []domain.Role) (domain.User, error) {
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

	usr = domain.NewUser(id, email, username, firstName, lastName, roles)

	userKey, err := s.enc.GenerateUserEncryptionKey()
	if err != nil {
		return domain.User{}, err
	}

	encryptedUserKey, err := s.enc.EncryptUserKey(userKey)
	if err != nil {
		return domain.User{}, err
	}

	ent, err := entities.ToEntity(usr, s.enc, encryptedUserKey, userKey)
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

	usr, err := entities.FromEntity(ent, s.enc)

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

	usr, err := entities.FromEntity(ent, s.enc)

	if err != nil {
		return domain.User{}, err
	}

	return usr, nil
}

func (s *service) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	ent, err := s.repo.FindByUsername(ctx, username)

	if err != nil {
		return domain.User{}, err
	}

	usr, err := entities.FromEntity(ent, s.enc)

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

func (s *service) UpdateIntroduction(ctx context.Context, id domain.UserId, introduction string) error {
	encryptedUserKey, err := s.repo.GetEncryptedUserKey(ctx, id.UUID)
	if err != nil {
		return err
	}

	userKey, err := s.enc.DecryptUserKey(encryptedUserKey)
	if err != nil {
		return err
	}

	encrypted, err := s.enc.Encrypt([]byte(introduction), userKey)
	if err != nil {
		return err
	}

	return s.repo.UpdateIntroduction(ctx, id.UUID, encrypted)
}

func (s *service) UpdateAbout(ctx context.Context, id domain.UserId, about string) error {
	encryptedUserKey, err := s.repo.GetEncryptedUserKey(ctx, id.UUID)
	if err != nil {
		return err
	}

	userKey, err := s.enc.DecryptUserKey(encryptedUserKey)
	if err != nil {
		return err
	}

	encrypted, err := s.enc.Encrypt([]byte(about), userKey)
	if err != nil {
		return err
	}

	return s.repo.UpdateAbout(ctx, id.UUID, encrypted)
}

func (s *service) UpdateImage(ctx context.Context, id domain.UserId, image string) error {
	return s.repo.UpdateImage(ctx, id.UUID, image)
}

func (s *service) UpsertAddress(ctx context.Context, userID domain.UserId, street string, houseNumber string, postalCode string, city string, country string) (domain.Address, error) {
	addrID, err := domain.NewAddressId()
	if err != nil {
		return domain.Address{}, err
	}

	encryptedUserKey, err := s.repo.GetEncryptedUserKey(ctx, userID.UUID)
	if err != nil {
		return domain.Address{}, err
	}

	userKey, err := s.enc.DecryptUserKey(encryptedUserKey)
	if err != nil {
		return domain.Address{}, err
	}

	addr := domain.NewAddress(addrID, userID, street, houseNumber, postalCode, city, country)

	ent, err := entities.AddressToEntity(addr, userKey, s.enc)
	if err != nil {
		return domain.Address{}, err
	}

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

	encryptedUserKey, err := s.repo.GetEncryptedUserKey(ctx, userID.UUID)
	if err != nil {
		return domain.Address{}, err
	}

	userKey, err := s.enc.DecryptUserKey(encryptedUserKey)
	if err != nil {
		return domain.Address{}, err
	}

	return entities.AddressFromEntity(ent, userKey, s.enc)
}

func (s *service) UpsertDevice(ctx context.Context, userID domain.UserId, deviceToken string, platform string) error {
	id, err := uuid.NewV4()
	if err != nil {
		return err
	}

	now := time.Now()

	ent := entities.UserDeviceEntity{
		ID:          id,
		UserID:      userID.UUID,
		DeviceToken: deviceToken,
		Platform:    platform,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	return s.repo.UpsertDevice(ctx, ent)
}

func (s *service) DeleteDevice(ctx context.Context, deviceToken string) error {
	return s.repo.DeleteDevice(ctx, deviceToken)
}
