package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/google/uuid"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/infrastructure/entities"
)

func (r *repository) Save(ctx context.Context, e entities.UserEntity) error {
	query := `
		INSERT INTO users (id, email, roles)
		VALUES ($1, $2, $3)
		ON CONFLICT (id) DO UPDATE
		SET email = EXCLUDED.email,
		    roles = EXCLUDED.roles,
		    updated_at = now()
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		e.ID,
		e.Email,
		e.Roles,
	)
	return err
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (entities.UserEntity, error) {
	var e entities.UserEntity

	query := `
		SELECT id, email, roles, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	if err := r.db.
		QueryRowContext(ctx, query, id).
		Scan(&e.ID, &e.Email, &e.Roles, &e.CreatedAt, &e.UpdatedAt); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return entities.UserEntity{}, ErrNotFound
		}
		return entities.UserEntity{}, err
	}

	return e, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (entities.UserEntity, error) {
	var e entities.UserEntity

	query := `
		SELECT id, email, roles, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	if err := r.db.
		QueryRowContext(ctx, query, email).
		Scan(&e.ID, &e.Email, &e.Roles, &e.CreatedAt, &e.UpdatedAt); err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return entities.UserEntity{}, ErrNotFound
		}
		return entities.UserEntity{}, err
	}

	return e, nil
}

func (r *repository) AddSupportMember(ctx context.Context, veteranID, supportID uuid.UUID) error {
	query := `
		INSERT INTO user_supports (veteran_id, support_id)
		VALUES ($1, $2)
		ON CONFLICT (veteran_id, support_id) DO NOTHING
	`

	_, err := r.db.ExecContext(ctx, query, veteranID, supportID)
	return err
}

func (r *repository) ListSupportMembers(ctx context.Context, veteranID uuid.UUID) ([]entities.UserEntity, error) {
	query := `
		SELECT u.id, u.email, u.roles, u.created_at, u.updated_at
		FROM users u
		JOIN user_supports s ON u.id = s.support_id
		WHERE s.veteran_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, veteranID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	var out []entities.UserEntity
	for rows.Next() {
		var e entities.UserEntity
		if err := rows.Scan(&e.ID, &e.Email, &e.Roles, &e.CreatedAt, &e.UpdatedAt); err != nil {
			return nil, err
		}
		out = append(out, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return out, nil
}
