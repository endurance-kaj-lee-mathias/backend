package infrastructure

import (
	"context"
	"log/slog"
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/support/infrastructure/entities"
)

func (r *repository) Create(ctx context.Context, veteranID, memberId uuid.UUID) (entities.MemberEntity, error) {
	var ent entities.MemberEntity

	query := `
		INSERT INTO user_supports (veteran_id, member_id, created_at)
		VALUES ($1, $2, $3)
		RETURNING veteran_id, member_id, created_at
		ON CONFLICT (veteran_id, member_id) DO NOTHING
	`

	err := r.db.QueryRow(query, veteranID, memberId, time.Now().UTC()).Scan(&ent.ID, &ent.Veteran, &ent.CreatedAt)

	if err != nil {
		return entities.MemberEntity{}, err
	}

	return ent, nil
}

func (r *repository) ReadAll(ctx context.Context, id uuid.UUID) ([]entities.MemberEntity, error) {
	query := `
		SELECT u.id, u.email, u.roles, u.created_at, u.updated_at
		FROM users u
		JOIN user_supports s ON u.id = s.member_id
		WHERE s.veteran_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, id)

	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	var ents []entities.MemberEntity

	for rows.Next() {
		var ent entities.MemberEntity

		err := rows.Scan(
			&ent.ID, &ent.Veteran, &ent.Email, &ent.CreatedAt, &ent.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		ents = append(ents, ent)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return ents, nil
}
