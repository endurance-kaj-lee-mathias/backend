package entities

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/encryption"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
)

type UserEntity struct {
	ID                   uuid.UUID `db:"id"`
	EmailHash            string    `db:"email_hash"`
	UsernameHash         string    `db:"username_hash"`
	PhoneNumberHash      *string   `db:"phone_number_hash"`
	EncryptedEmail       []byte    `db:"encrypted_email"`
	EncryptedUsername    []byte    `db:"encrypted_username"`
	EncryptedFirstName   []byte    `db:"encrypted_first_name"`
	EncryptedLastName    []byte    `db:"encrypted_last_name"`
	EncryptedPhoneNumber []byte    `db:"encrypted_phone_number"`
	EncryptedRoles       []byte    `db:"encrypted_roles"`
	EncryptedUserKey     []byte    `db:"encrypted_user_key"`
	KeyVersion           int       `db:"key_version"`
	CreatedAt            time.Time `db:"created_at"`
	UpdatedAt            time.Time `db:"updated_at"`
}

func FromEntity(ent UserEntity, enc encryption.Service) (domain.User, error) {
	userKey, err := enc.DecryptUserKey(ent.EncryptedUserKey)
	if err != nil {
		return domain.User{}, err
	}

	emailBytes, err := enc.Decrypt(ent.EncryptedEmail, userKey)
	if err != nil {
		return domain.User{}, err
	}

	usernameBytes, err := enc.Decrypt(ent.EncryptedUsername, userKey)
	if err != nil {
		return domain.User{}, err
	}

	firstNameBytes, err := enc.Decrypt(ent.EncryptedFirstName, userKey)
	if err != nil {
		return domain.User{}, err
	}

	lastNameBytes, err := enc.Decrypt(ent.EncryptedLastName, userKey)
	if err != nil {
		return domain.User{}, err
	}

	var phoneNumber *string
	if len(ent.EncryptedPhoneNumber) > 0 {
		phoneBytes, err := enc.Decrypt(ent.EncryptedPhoneNumber, userKey)
		if err != nil {
			return domain.User{}, err
		}
		phone := string(phoneBytes)
		phoneNumber = &phone
	}

	rolesBytes, err := enc.Decrypt(ent.EncryptedRoles, userKey)
	if err != nil {
		return domain.User{}, err
	}

	roles := make([]domain.Role, 0)
	if len(rolesBytes) > 0 {
		if err := json.Unmarshal(rolesBytes, &roles); err != nil {
			return domain.User{}, InvalidRoles
		}
	}

	id, err := domain.ParseId(ent.ID.String())
	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		ID:          id,
		Email:       string(emailBytes),
		Username:    string(usernameBytes),
		FirstName:   string(firstNameBytes),
		LastName:    string(lastNameBytes),
		PhoneNumber: phoneNumber,
		Roles:       roles,
		CreatedAt:   ent.CreatedAt,
		UpdatedAt:   ent.UpdatedAt,
	}, nil
}

func ToEntity(usr domain.User, enc encryption.Service, encryptedUserKey []byte, userKey []byte) (UserEntity, error) {
	encEmail, err := enc.Encrypt([]byte(usr.Email), userKey)
	if err != nil {
		return UserEntity{}, err
	}

	encUsername, err := enc.Encrypt([]byte(usr.Username), userKey)
	if err != nil {
		return UserEntity{}, err
	}

	encFirstName, err := enc.Encrypt([]byte(usr.FirstName), userKey)
	if err != nil {
		return UserEntity{}, err
	}

	encLastName, err := enc.Encrypt([]byte(usr.LastName), userKey)
	if err != nil {
		return UserEntity{}, err
	}

	var encPhoneNumber []byte
	var phoneNumberHash *string
	if usr.PhoneNumber != nil {
		encPhone, err := enc.Encrypt([]byte(*usr.PhoneNumber), userKey)
		if err != nil {
			return UserEntity{}, err
		}
		encPhoneNumber = encPhone
		hash := enc.Hash(*usr.PhoneNumber)
		phoneNumberHash = &hash
	}

	rolesJSON, err := json.Marshal(usr.Roles)
	if err != nil {
		return UserEntity{}, InvalidRoles
	}

	encRoles, err := enc.Encrypt(rolesJSON, userKey)
	if err != nil {
		return UserEntity{}, err
	}

	return UserEntity{
		ID:                   usr.ID.UUID,
		EmailHash:            enc.Hash(usr.Email),
		UsernameHash:         enc.Hash(usr.Username),
		PhoneNumberHash:      phoneNumberHash,
		EncryptedEmail:       encEmail,
		EncryptedUsername:    encUsername,
		EncryptedFirstName:   encFirstName,
		EncryptedLastName:    encLastName,
		EncryptedPhoneNumber: encPhoneNumber,
		EncryptedRoles:       encRoles,
		EncryptedUserKey:     encryptedUserKey,
		KeyVersion:           1,
		CreatedAt:            usr.CreatedAt,
		UpdatedAt:            usr.UpdatedAt,
	}, nil
}
