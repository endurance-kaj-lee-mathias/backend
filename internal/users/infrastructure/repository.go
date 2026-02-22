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
		INSERT INTO users (id, email, first_name, last_name, phone_number, roles, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (id) DO UPDATE
		SET email = EXCLUDED.email,
			first_name = EXCLUDED.first_name,
			last_name = EXCLUDED.last_name,
		    roles = EXCLUDED.roles,
		    updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		ent.ID,
		ent.Email,
		ent.FirstName,
		ent.LastName,
		ent.PhoneNumber,
		ent.Roles,
		ent.CreatedAt,
		ent.UpdatedAt,
	)

	return err
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (entities.UserEntity, error) {
	var e entities.UserEntity

	query := `
		SELECT id, email, first_name, last_name, phone_number, roles, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	if err := r.db.
		QueryRowContext(ctx, query, id).
		Scan(&e.ID, &e.Email, &e.FirstName, &e.LastName, &e.PhoneNumber, &e.Roles, &e.CreatedAt, &e.UpdatedAt); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return entities.UserEntity{}, NotFound
		}

		return entities.UserEntity{}, err
	}

	return e, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (entities.UserEntity, error) {
	var ent entities.UserEntity

	query := `
		SELECT id, email, first_name, last_name, phone_number, roles, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	if err := r.db.
		QueryRowContext(ctx, query, email).
		Scan(&ent.ID, &ent.Email, &ent.FirstName, &ent.LastName, &ent.PhoneNumber, &ent.Roles, &ent.CreatedAt, &ent.UpdatedAt); err != nil {

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
	query := `UPDATE users SET phone_number = $1, updated_at = NOW() WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, phoneNumber, id)
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
		INSERT INTO user_addresses (id, user_id, street, house_number, postal_code, city, country, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id) DO UPDATE
		SET street       = EXCLUDED.street,
		    house_number = EXCLUDED.house_number,
		    postal_code  = EXCLUDED.postal_code,
		    city         = EXCLUDED.city,
		    country      = EXCLUDED.country
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		ent.ID,
		ent.UserID,
		ent.Street,
		ent.HouseNumber,
		ent.PostalCode,
		ent.City,
		ent.Country,
		ent.CreatedAt,
	)

	return err
}

func (r *repository) FindAddressByUserID(ctx context.Context, userID uuid.UUID) (entities.AddressEntity, error) {
	var ent entities.AddressEntity

	query := `
		SELECT id, user_id, street, house_number, postal_code, city, country, created_at
		FROM user_addresses
		WHERE user_id = $1
	`

	if err := r.db.
		QueryRowContext(ctx, query, userID).
		Scan(&ent.ID, &ent.UserID, &ent.Street, &ent.HouseNumber, &ent.PostalCode, &ent.City, &ent.Country, &ent.CreatedAt); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return entities.AddressEntity{}, AddressNotFound
		}

		return entities.AddressEntity{}, err
	}

	return ent, nil
}
