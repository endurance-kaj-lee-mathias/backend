package infrastructure

import (
	"context"
	"log/slog"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/infrastructure/entities"
)

func (r *repository) GetEventsByUserID(ctx context.Context, userID uuid.UUID) ([]entities.CalendarEventEntity, error) {
	query := `
		SELECT a.id, s.start_time, s.end_time, a.updated_at,
			u.encrypted_first_name, u.encrypted_last_name, u.encrypted_user_key
		FROM appointments a
		JOIN availability_slots s ON s.id = a.slot_id
		JOIN users u ON u.id = s.provider_id
		WHERE a.veteran_id = $1
		AND a.status != 'CANCELLED'
		ORDER BY s.start_time ASC
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

	var result []entities.CalendarEventEntity

	for rows.Next() {
		var ent entities.CalendarEventEntity
		if err := rows.Scan(
			&ent.ID, &ent.StartTime, &ent.EndTime, &ent.UpdatedAt,
			&ent.ProviderEncryptedFirstName, &ent.ProviderEncryptedLastName, &ent.ProviderEncryptedUserKey,
		); err != nil {
			return nil, err
		}
		result = append(result, ent)
	}

	return result, rows.Err()
}
