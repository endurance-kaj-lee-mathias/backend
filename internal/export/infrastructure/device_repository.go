package infrastructure

import (
	"context"
	"log/slog"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/infrastructure/entities"
)

func (r *repository) GetDevices(ctx context.Context, userID uuid.UUID) ([]entities.DeviceExportEntity, error) {
	query := `
		SELECT device_token, platform, created_at
		FROM user_devices
		WHERE user_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	var result []entities.DeviceExportEntity
	for rows.Next() {
		var ent entities.DeviceExportEntity
		if err := rows.Scan(&ent.DeviceToken, &ent.Platform, &ent.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, ent)
	}

	return result, rows.Err()
}
