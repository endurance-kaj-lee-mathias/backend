package infrastructure

import (
	"context"
	"log/slog"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/infrastructure/entities"
)

func (r *repository) GetAuthorizationRulesAsOwner(ctx context.Context, userID uuid.UUID) ([]entities.AuthorizationRuleExportEntity, error) {
	return r.queryAuthorizationRules(ctx, `
		SELECT id, owner_id, viewer_id, resource, effect, created_at
		FROM authorization_rules
		WHERE owner_id = $1
		ORDER BY created_at DESC
	`, userID)
}

func (r *repository) GetAuthorizationRulesAsViewer(ctx context.Context, userID uuid.UUID) ([]entities.AuthorizationRuleExportEntity, error) {
	return r.queryAuthorizationRules(ctx, `
		SELECT id, owner_id, viewer_id, resource, effect, created_at
		FROM authorization_rules
		WHERE viewer_id = $1
		ORDER BY created_at DESC
	`, userID)
}

func (r *repository) queryAuthorizationRules(ctx context.Context, query string, userID uuid.UUID) ([]entities.AuthorizationRuleExportEntity, error) {
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	var result []entities.AuthorizationRuleExportEntity
	for rows.Next() {
		var ent entities.AuthorizationRuleExportEntity
		if err := rows.Scan(&ent.ID, &ent.OwnerID, &ent.ViewerID, &ent.Resource, &ent.Effect, &ent.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, ent)
	}

	return result, rows.Err()
}
