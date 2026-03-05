package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/infrastructure/entities"
)

func (r *repository) Create(ctx context.Context, ent entities.UserEntity) error {
	query := `
		INSERT INTO users (id, email_hash, username_hash, phone_number_hash, role_hash, encrypted_email, encrypted_username, encrypted_first_name, encrypted_last_name, encrypted_phone_number, encrypted_roles, encrypted_about, encrypted_introduction, image, is_private, encrypted_user_key, key_version, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		ON CONFLICT (id) DO UPDATE
		SET email_hash             = EXCLUDED.email_hash,
			username_hash          = EXCLUDED.username_hash,
			phone_number_hash      = EXCLUDED.phone_number_hash,
			role_hash              = EXCLUDED.role_hash,
			encrypted_email        = EXCLUDED.encrypted_email,
			encrypted_username     = EXCLUDED.encrypted_username,
			encrypted_first_name   = EXCLUDED.encrypted_first_name,
			encrypted_last_name    = EXCLUDED.encrypted_last_name,
			encrypted_phone_number = EXCLUDED.encrypted_phone_number,
			encrypted_roles        = EXCLUDED.encrypted_roles,
			encrypted_about        = EXCLUDED.encrypted_about,
			encrypted_introduction = EXCLUDED.encrypted_introduction,
			image                  = EXCLUDED.image,
			encrypted_user_key     = EXCLUDED.encrypted_user_key,
			key_version            = EXCLUDED.key_version,
			updated_at             = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		ent.ID,
		ent.EmailHash,
		ent.UsernameHash,
		ent.PhoneNumberHash,
		ent.RoleHash,
		ent.EncryptedEmail,
		ent.EncryptedUsername,
		ent.EncryptedFirstName,
		ent.EncryptedLastName,
		ent.EncryptedPhoneNumber,
		ent.EncryptedRoles,
		ent.EncryptedAbout,
		ent.EncryptedIntroduction,
		ent.Image,
		ent.IsPrivate,
		ent.EncryptedUserKey,
		ent.KeyVersion,
		ent.CreatedAt,
		ent.UpdatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && strings.Contains(pgErr.ConstraintName, "username") {
			return UsernameAlreadyExists
		}
		return err
	}

	return nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (entities.UserEntity, error) {
	var e entities.UserEntity

	query := `
		SELECT id, email_hash, username_hash, phone_number_hash, encrypted_email, encrypted_username, encrypted_first_name, encrypted_last_name, encrypted_phone_number, encrypted_roles, encrypted_about, encrypted_introduction, image, is_private, encrypted_user_key, key_version, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	if err := r.db.
		QueryRowContext(ctx, query, id).
		Scan(&e.ID, &e.EmailHash, &e.UsernameHash, &e.PhoneNumberHash, &e.EncryptedEmail, &e.EncryptedUsername, &e.EncryptedFirstName, &e.EncryptedLastName, &e.EncryptedPhoneNumber, &e.EncryptedRoles, &e.EncryptedAbout, &e.EncryptedIntroduction, &e.Image, &e.IsPrivate, &e.EncryptedUserKey, &e.KeyVersion, &e.CreatedAt, &e.UpdatedAt); err != nil {

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
		SELECT id, email_hash, username_hash, phone_number_hash, encrypted_email, encrypted_username, encrypted_first_name, encrypted_last_name, encrypted_phone_number, encrypted_roles, encrypted_about, encrypted_introduction, image, is_private, encrypted_user_key, key_version, created_at, updated_at
		FROM users
		WHERE email_hash = $1
	`

	if err := r.db.
		QueryRowContext(ctx, query, emailHash).
		Scan(&ent.ID, &ent.EmailHash, &ent.UsernameHash, &ent.PhoneNumberHash, &ent.EncryptedEmail, &ent.EncryptedUsername, &ent.EncryptedFirstName, &ent.EncryptedLastName, &ent.EncryptedPhoneNumber, &ent.EncryptedRoles, &ent.EncryptedAbout, &ent.EncryptedIntroduction, &ent.Image, &ent.IsPrivate, &ent.EncryptedUserKey, &ent.KeyVersion, &ent.CreatedAt, &ent.UpdatedAt); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return entities.UserEntity{}, NotFound
		}

		return entities.UserEntity{}, err
	}

	return ent, nil
}

func (r *repository) FindByUsername(ctx context.Context, username string) (entities.UserEntity, error) {
	var ent entities.UserEntity

	usernameHash := r.enc.Hash(username)

	query := `
		SELECT id, email_hash, username_hash, phone_number_hash, encrypted_email, encrypted_username, encrypted_first_name, encrypted_last_name, encrypted_phone_number, encrypted_roles, encrypted_about, encrypted_introduction, image, is_private, encrypted_user_key, key_version, created_at, updated_at
		FROM users
		WHERE username_hash = $1
	`

	if err := r.db.
		QueryRowContext(ctx, query, usernameHash).
		Scan(&ent.ID, &ent.EmailHash, &ent.UsernameHash, &ent.PhoneNumberHash, &ent.EncryptedEmail, &ent.EncryptedUsername, &ent.EncryptedFirstName, &ent.EncryptedLastName, &ent.EncryptedPhoneNumber, &ent.EncryptedRoles, &ent.EncryptedAbout, &ent.EncryptedIntroduction, &ent.Image, &ent.IsPrivate, &ent.EncryptedUserKey, &ent.KeyVersion, &ent.CreatedAt, &ent.UpdatedAt); err != nil {

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

func (r *repository) UpdateFirstName(ctx context.Context, id uuid.UUID, encrypted []byte) error {
	query := `UPDATE users SET encrypted_first_name = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, encrypted, id)
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

func (r *repository) UpdateLastName(ctx context.Context, id uuid.UUID, encrypted []byte) error {
	query := `UPDATE users SET encrypted_last_name = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, encrypted, id)
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

func (r *repository) UpdateIntroduction(ctx context.Context, id uuid.UUID, encrypted []byte) error {
	query := `UPDATE users SET encrypted_introduction = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, encrypted, id)
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

func (r *repository) UpdateAbout(ctx context.Context, id uuid.UUID, encrypted []byte) error {
	query := `UPDATE users SET encrypted_about = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, encrypted, id)
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

func (r *repository) UpdateImage(ctx context.Context, id uuid.UUID, image string) error {
	query := `UPDATE users SET image = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, image, id)
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

func (r *repository) UpdatePrivacy(ctx context.Context, id uuid.UUID, isPrivate bool) error {
	query := `UPDATE users SET is_private = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, isPrivate, id)
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
		INSERT INTO user_addresses (id, user_id, encrypted_street, encrypted_locality, encrypted_region, encrypted_postal_code, encrypted_country, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id) DO UPDATE
		SET encrypted_street       = EXCLUDED.encrypted_street,
			encrypted_locality     = EXCLUDED.encrypted_locality,
			encrypted_region       = EXCLUDED.encrypted_region,
			encrypted_postal_code  = EXCLUDED.encrypted_postal_code,
			encrypted_country      = EXCLUDED.encrypted_country
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		ent.ID,
		ent.UserID,
		ent.EncryptedStreet,
		ent.EncryptedLocality,
		ent.EncryptedRegion,
		ent.EncryptedPostalCode,
		ent.EncryptedCountry,
		ent.CreatedAt,
	)

	return err
}

func (r *repository) FindAddressByUserID(ctx context.Context, userID uuid.UUID) (entities.AddressEntity, error) {
	var ent entities.AddressEntity

	query := `
		SELECT id, user_id, encrypted_street, encrypted_locality, encrypted_region, encrypted_postal_code, encrypted_country, created_at
		FROM user_addresses
		WHERE user_id = $1
	`

	if err := r.db.
		QueryRowContext(ctx, query, userID).
		Scan(&ent.ID, &ent.UserID, &ent.EncryptedStreet, &ent.EncryptedLocality, &ent.EncryptedRegion, &ent.EncryptedPostalCode, &ent.EncryptedCountry, &ent.CreatedAt); err != nil {

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

func (r *repository) UpsertDevice(ctx context.Context, ent entities.UserDeviceEntity) error {
	query := `
		INSERT INTO user_devices (id, user_id, device_token, platform, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (device_token) DO UPDATE
		SET updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(ctx, query, ent.ID, ent.UserID, ent.DeviceToken, ent.Platform, ent.CreatedAt, ent.UpdatedAt)
	return err
}

func (r *repository) DeleteDevice(ctx context.Context, deviceToken string) error {
	query := `DELETE FROM user_devices WHERE device_token = $1`

	result, err := r.db.ExecContext(ctx, query, deviceToken)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return DeviceNotFound
	}

	return nil
}

func (r *repository) FindDeviceTokensByUserID(ctx context.Context, userID uuid.UUID) ([]string, error) {
	query := `SELECT device_token FROM user_devices WHERE user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	return tokens, rows.Err()
}
