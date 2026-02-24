package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/infrastructure/entities"
)

func (r *repository) Create(ctx context.Context, ent entities.UserEntity) error {
	query := `
		INSERT INTO users (id, email_hash, phone_number_hash, encrypted_email, encrypted_first_name, encrypted_last_name, encrypted_phone_number, encrypted_roles, encrypted_user_key, key_version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (id) DO UPDATE
		SET email_hash             = EXCLUDED.email_hash,
			phone_number_hash      = EXCLUDED.phone_number_hash,
			encrypted_email        = EXCLUDED.encrypted_email,
			encrypted_first_name   = EXCLUDED.encrypted_first_name,
			encrypted_last_name    = EXCLUDED.encrypted_last_name,
			encrypted_phone_number = EXCLUDED.encrypted_phone_number,
			encrypted_roles        = EXCLUDED.encrypted_roles,
			encrypted_user_key     = EXCLUDED.encrypted_user_key,
			key_version            = EXCLUDED.key_version,
			updated_at             = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		ent.ID,
		ent.EmailHash,
		ent.PhoneNumberHash,
		ent.EncryptedEmail,
		ent.EncryptedFirstName,
		ent.EncryptedLastName,
		ent.EncryptedPhoneNumber,
		ent.EncryptedRoles,
		ent.EncryptedUserKey,
		ent.KeyVersion,
		ent.CreatedAt,
		ent.UpdatedAt,
	)

	return err
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (entities.UserEntity, error) {
	var e entities.UserEntity

	query := `
		SELECT id, email_hash, phone_number_hash, encrypted_email, encrypted_first_name, encrypted_last_name, encrypted_phone_number, encrypted_roles, encrypted_user_key, key_version, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	if err := r.db.
		QueryRowContext(ctx, query, id).
		Scan(&e.ID, &e.EmailHash, &e.PhoneNumberHash, &e.EncryptedEmail, &e.EncryptedFirstName, &e.EncryptedLastName, &e.EncryptedPhoneNumber, &e.EncryptedRoles, &e.EncryptedUserKey, &e.KeyVersion, &e.CreatedAt, &e.UpdatedAt); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return entities.UserEntity{}, NotFound
		}

		return entities.UserEntity{}, err
	}

	return e, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (entities.UserEntity, error) {
	var ent entities.UserEntity

	emailHash := r.enc.Hash(email)

	query := `
		SELECT id, email_hash, phone_number_hash, encrypted_email, encrypted_first_name, encrypted_last_name, encrypted_phone_number, encrypted_roles, encrypted_user_key, key_version, created_at, updated_at
		FROM users
		WHERE email_hash = $1
	`

	if err := r.db.
		QueryRowContext(ctx, query, emailHash).
		Scan(&ent.ID, &ent.EmailHash, &ent.PhoneNumberHash, &ent.EncryptedEmail, &ent.EncryptedFirstName, &ent.EncryptedLastName, &ent.EncryptedPhoneNumber, &ent.EncryptedRoles, &ent.EncryptedUserKey, &ent.KeyVersion, &ent.CreatedAt, &ent.UpdatedAt); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return entities.UserEntity{}, NotFound
		}

		return entities.UserEntity{}, err
	}

	return ent, nil
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return NotFound
	}

	return nil
}

func (r *repository) UpdatePhoneNumber(ctx context.Context, id uuid.UUID, phoneNumber *string) error {
	encryptedUserKey, err := r.GetEncryptedUserKey(ctx, id)
	if err != nil {
		return err
	}

	userKey, err := r.enc.DecryptUserKey(encryptedUserKey)
	if err != nil {
		return err
	}

	var encPhoneNumber []byte
	var phoneNumberHash *string
	if phoneNumber != nil {
		encPhone, err := r.enc.Encrypt([]byte(*phoneNumber), userKey)
		if err != nil {
			return err
		}
		encPhoneNumber = encPhone
		hash := r.enc.Hash(*phoneNumber)
		phoneNumberHash = &hash
	}

	query := `UPDATE users SET encrypted_phone_number = $1, phone_number_hash = $2, updated_at = NOW() WHERE id = $3`

	result, err := r.db.ExecContext(ctx, query, encPhoneNumber, phoneNumberHash, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return NotFound
	}

	return nil
}

func (r *repository) InsertAddress(ctx context.Context, ent entities.AddressEntity) error {
	query := `
		INSERT INTO user_addresses (id, user_id, encrypted_street, encrypted_house_number, encrypted_postal_code, encrypted_city, encrypted_country, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id) DO UPDATE
		SET encrypted_street       = EXCLUDED.encrypted_street,
			encrypted_house_number = EXCLUDED.encrypted_house_number,
			encrypted_postal_code  = EXCLUDED.encrypted_postal_code,
			encrypted_city         = EXCLUDED.encrypted_city,
			encrypted_country      = EXCLUDED.encrypted_country
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		ent.ID,
		ent.UserID,
		ent.EncryptedStreet,
		ent.EncryptedHouseNumber,
		ent.EncryptedPostalCode,
		ent.EncryptedCity,
		ent.EncryptedCountry,
		ent.CreatedAt,
	)

	return err
}

func (r *repository) FindAddressByUserID(ctx context.Context, userID uuid.UUID) (entities.AddressEntity, error) {
	var ent entities.AddressEntity

	query := `
		SELECT id, user_id, encrypted_street, encrypted_house_number, encrypted_postal_code, encrypted_city, encrypted_country, created_at
		FROM user_addresses
		WHERE user_id = $1
	`

	if err := r.db.
		QueryRowContext(ctx, query, userID).
		Scan(&ent.ID, &ent.UserID, &ent.EncryptedStreet, &ent.EncryptedHouseNumber, &ent.EncryptedPostalCode, &ent.EncryptedCity, &ent.EncryptedCountry, &ent.CreatedAt); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return entities.AddressEntity{}, AddressNotFound
		}

		return entities.AddressEntity{}, err
	}

	return ent, nil
}

func (r *repository) GetEncryptedUserKey(ctx context.Context, userID uuid.UUID) ([]byte, error) {
	var encryptedUserKey []byte

	query := `SELECT encrypted_user_key FROM users WHERE id = $1`

	if err := r.db.QueryRowContext(ctx, query, userID).Scan(&encryptedUserKey); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, NotFound
		}
		return nil, err
	}

	return encryptedUserKey, nil
}
