package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/authorization/domain"
)

func (r *repository) Create(ctx context.Context, rule domain.Rule) error {
	query := `
		INSERT INTO authorization_rules (id, owner_id, viewer_id, resource, effect, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (owner_id, viewer_id, resource) DO UPDATE
		SET effect = EXCLUDED.effect
	`

	_, err := r.db.ExecContext(ctx, query, rule.ID, rule.OwnerID, rule.ViewerID, rule.Resource, rule.Effect, rule.CreatedAt)
	return err
}

func (r *repository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM authorization_rules WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return RuleNotFound
	}

	return nil
}

func (r *repository) FindByID(ctx context.Context, id uuid.UUID) (domain.Rule, error) {
	query := `
		SELECT id, owner_id, viewer_id, resource, effect, created_at
		FROM authorization_rules
		WHERE id = $1
	`

	var rule domain.Rule
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rule.ID, &rule.OwnerID, &rule.ViewerID, &rule.Resource, &rule.Effect, &rule.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Rule{}, RuleNotFound
		}
		return domain.Rule{}, err
	}

	return rule, nil
}

func (r *repository) FindByOwner(ctx context.Context, ownerID uuid.UUID) ([]domain.Rule, error) {
	query := `
		SELECT id, owner_id, viewer_id, resource, effect, created_at
		FROM authorization_rules
		WHERE owner_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	var rules []domain.Rule
	for rows.Next() {
		var rule domain.Rule
		if err := rows.Scan(&rule.ID, &rule.OwnerID, &rule.ViewerID, &rule.Resource, &rule.Effect, &rule.CreatedAt); err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	return rules, rows.Err()
}

func (r *repository) FindByOwnerAndViewer(ctx context.Context, ownerID uuid.UUID, viewerID uuid.UUID) ([]domain.Rule, error) {
	query := `
		SELECT id, owner_id, viewer_id, resource, effect, created_at
		FROM authorization_rules
		WHERE owner_id = $1 AND viewer_id = $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, ownerID, viewerID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	var rules []domain.Rule
	for rows.Next() {
		var rule domain.Rule
		if err := rows.Scan(&rule.ID, &rule.OwnerID, &rule.ViewerID, &rule.Resource, &rule.Effect, &rule.CreatedAt); err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	return rules, rows.Err()
}

func (r *repository) FindRule(ctx context.Context, ownerID uuid.UUID, viewerID uuid.UUID, resource string) (*domain.Rule, error) {
	query := `
		SELECT id, owner_id, viewer_id, resource, effect, created_at
		FROM authorization_rules
		WHERE owner_id = $1 AND viewer_id = $2 AND resource IN ($3, '*')
		ORDER BY CASE WHEN resource = $3 THEN 0 ELSE 1 END
		LIMIT 1
	`

	var rule domain.Rule
	err := r.db.QueryRowContext(ctx, query, ownerID, viewerID, resource).Scan(
		&rule.ID, &rule.OwnerID, &rule.ViewerID, &rule.Resource, &rule.Effect, &rule.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &rule, nil
}

func (r *repository) DeleteByOwnerAndViewer(ctx context.Context, ownerID uuid.UUID, viewerID uuid.UUID) error {
	query := `DELETE FROM authorization_rules WHERE owner_id = $1 AND viewer_id = $2`
	_, err := r.db.ExecContext(ctx, query, ownerID, viewerID)
	return err
}

func (r *repository) GetPrivacy(ctx context.Context, userID uuid.UUID) (bool, error) {
	var isPrivate bool

	query := `SELECT is_private FROM users WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&isPrivate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return isPrivate, nil
}

func (r *repository) HasSupportRelationship(ctx context.Context, userA uuid.UUID, userB uuid.UUID) (bool, error) {
	query := `
		SELECT 1 FROM user_supports
		WHERE (veteran_id = $1 AND support_id = $2) OR (veteran_id = $2 AND support_id = $1)
		LIMIT 1
	`

	var exists int
	err := r.db.QueryRowContext(ctx, query, userA, userB).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
