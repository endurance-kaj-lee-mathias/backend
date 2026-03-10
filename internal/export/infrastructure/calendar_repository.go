package infrastructure

import (
	"context"
	"log/slog"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/infrastructure/entities"
)

func (r *repository) GetAppointmentsWithSlots(ctx context.Context, userID uuid.UUID) ([]entities.AppointmentExportEntity, error) {
	query := `
		SELECT a.id, a.slot_id, a.veteran_id, s.provider_id, a.status,
			s.start_time, s.end_time, s.is_urgent,
			a.created_at, a.updated_at,
			u.encrypted_username, u.encrypted_user_key
		FROM appointments a
		JOIN availability_slots s ON s.id = a.slot_id
		JOIN users u ON u.id = s.provider_id
		WHERE a.veteran_id = $1
		ORDER BY s.start_time DESC
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

	var result []entities.AppointmentExportEntity
	for rows.Next() {
		var ent entities.AppointmentExportEntity
		if err := rows.Scan(
			&ent.ID, &ent.SlotID, &ent.VeteranID, &ent.ProviderID, &ent.Status,
			&ent.StartTime, &ent.EndTime, &ent.IsUrgent,
			&ent.CreatedAt, &ent.UpdatedAt,
			&ent.ProviderEncryptedUsername, &ent.ProviderEncryptedUserKey,
		); err != nil {
			return nil, err
		}
		result = append(result, ent)
	}

	return result, rows.Err()
}

func (r *repository) GetSlotsAsProvider(ctx context.Context, providerID uuid.UUID) ([]entities.SlotExportEntity, error) {
	query := `
		SELECT id, start_time, end_time, is_urgent, is_booked, created_at, updated_at
		FROM availability_slots
		WHERE provider_id = $1
		ORDER BY start_time DESC
	`

	rows, err := r.db.QueryContext(ctx, query, providerID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	var result []entities.SlotExportEntity
	for rows.Next() {
		var ent entities.SlotExportEntity
		if err := rows.Scan(&ent.ID, &ent.StartTime, &ent.EndTime, &ent.IsUrgent, &ent.IsBooked, &ent.CreatedAt, &ent.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, ent)
	}

	return result, rows.Err()
}
