package infrastructure

import (
	"context"
	"log/slog"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/infrastructure/entities"
)

func (r *repository) GetEventsByUserID(ctx context.Context, userID uuid.UUID) ([]entities.CalendarEventEntity, error) {
	query := `
		SELECT a.id, s.provider_id, a.veteran_id, s.start_time, s.end_time, a.updated_at,
			p.encrypted_first_name, p.encrypted_last_name, p.encrypted_user_key,
			v.encrypted_first_name, v.encrypted_last_name, v.encrypted_user_key
		FROM appointments a
		JOIN availability_slots s ON s.id = a.slot_id
		JOIN users p ON p.id = s.provider_id
		JOIN users v ON v.id = a.veteran_id
		WHERE (a.veteran_id = $1 OR s.provider_id = $1)
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
			&ent.ID, &ent.ProviderID, &ent.VeteranID, &ent.StartTime, &ent.EndTime, &ent.UpdatedAt,
			&ent.ProviderEncryptedFirstName, &ent.ProviderEncryptedLastName, &ent.ProviderEncryptedUserKey,
			&ent.VeteranEncryptedFirstName, &ent.VeteranEncryptedLastName, &ent.VeteranEncryptedUserKey,
		); err != nil {
			return nil, err
		}
		result = append(result, ent)
	}

	return result, rows.Err()
}
