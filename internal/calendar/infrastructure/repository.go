package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/calendar/infrastructure/entities"
)

func (r *repository) CreateSlot(ctx context.Context, ent entities.SlotEntity) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			slog.Error("failed to rollback transaction", "error", err)
		}
	}()

	_, err = tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock(hashtext($1::text))`, ent.ProviderID)
	if err != nil {
		return err
	}

	var overlaps bool
	err = tx.QueryRowContext(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM availability_slots
			WHERE provider_id = $1 AND start_time < $3 AND end_time > $2
		)`,
		ent.ProviderID, ent.StartTime, ent.EndTime,
	).Scan(&overlaps)
	if err != nil {
		return err
	}
	if overlaps {
		return SlotOverlapDB
	}

	_, err = tx.ExecContext(ctx,
		`INSERT INTO availability_slots (id, provider_id, start_time, end_time, is_urgent, is_booked, series_id, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		ent.ID, ent.ProviderID, ent.StartTime, ent.EndTime, ent.IsUrgent, ent.IsBooked, ent.SeriesID, ent.CreatedAt, ent.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && strings.Contains(pgErr.ConstraintName, "no_overlap") {
			return SlotOverlapDB
		}
		return err
	}

	return tx.Commit()
}

func (r *repository) GetSlotsByRange(ctx context.Context, from, to time.Time, providerID *uuid.UUID) ([]entities.SlotEntity, error) {
	var rows *sql.Rows
	var err error

	if providerID != nil {
		rows, err = r.db.QueryContext(ctx,
			`SELECT s.id, s.provider_id, s.start_time, s.end_time, s.is_urgent, s.is_booked, s.series_id, s.created_at, s.updated_at, a.title
			 FROM availability_slots s
			 LEFT JOIN appointments a ON a.slot_id = s.id AND a.status != 'CANCELLED'
			 WHERE s.start_time >= $1 AND s.end_time <= $2 AND s.provider_id = $3
			 ORDER BY s.start_time`,
			from, to, *providerID,
		)
	} else {
		rows, err = r.db.QueryContext(ctx,
			`SELECT s.id, s.provider_id, s.start_time, s.end_time, s.is_urgent, s.is_booked, s.series_id, s.created_at, s.updated_at, a.title
			 FROM availability_slots s
			 LEFT JOIN appointments a ON a.slot_id = s.id AND a.status != 'CANCELLED'
			 WHERE s.start_time >= $1 AND s.end_time <= $2
			 ORDER BY s.start_time`,
			from, to,
		)
	}

	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	var ents []entities.SlotEntity
	for rows.Next() {
		var ent entities.SlotEntity
		if err := rows.Scan(
			&ent.ID, &ent.ProviderID, &ent.StartTime, &ent.EndTime,
			&ent.IsUrgent, &ent.IsBooked, &ent.SeriesID, &ent.CreatedAt, &ent.UpdatedAt, &ent.AppointmentTitle,
		); err != nil {
			return nil, err
		}
		ents = append(ents, ent)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return ents, nil
}

func (r *repository) GetSlotByID(ctx context.Context, id uuid.UUID) (entities.SlotEntity, error) {
	var ent entities.SlotEntity

	err := r.db.QueryRowContext(ctx,
		`SELECT s.id, s.provider_id, s.start_time, s.end_time, s.is_urgent, s.is_booked, s.series_id, s.created_at, s.updated_at, a.title
		 FROM availability_slots s
		 LEFT JOIN appointments a ON a.slot_id = s.id AND a.status != 'CANCELLED'
		 WHERE s.id = $1`,
		id,
	).Scan(&ent.ID, &ent.ProviderID, &ent.StartTime, &ent.EndTime, &ent.IsUrgent, &ent.IsBooked, &ent.SeriesID, &ent.CreatedAt, &ent.UpdatedAt, &ent.AppointmentTitle)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.SlotEntity{}, SlotNotFound
		}
		return entities.SlotEntity{}, err
	}

	return ent, nil
}

func (r *repository) GetSlotWithProvider(ctx context.Context, id uuid.UUID) (entities.SlotWithProviderEntity, error) {
	var ent entities.SlotWithProviderEntity

	err := r.db.QueryRowContext(ctx,
		`SELECT s.id, s.provider_id, s.start_time, s.end_time, s.is_urgent, s.is_booked, s.series_id , s.created_at, s.updated_at, a.title, u.encrypted_username, u.encrypted_first_name, u.encrypted_last_name, u.encrypted_user_key, u.image
		 FROM availability_slots s
		 LEFT JOIN appointments a ON a.slot_id = s.id AND a.status != 'CANCELLED'
		 JOIN users u ON u.id = s.provider_id
		 WHERE s.id = $1`,
		id,
	).Scan(&ent.ID, &ent.ProviderID, &ent.StartTime, &ent.EndTime, &ent.IsUrgent, &ent.IsBooked, &ent.SeriesID, &ent.CreatedAt, &ent.UpdatedAt, &ent.AppointmentTitle, &ent.ProviderUsernameEncrypted, &ent.ProviderFirstNameEncrypted, &ent.ProviderLastNameEncrypted, &ent.ProviderEncryptedUserKey, &ent.ProviderImage)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.SlotWithProviderEntity{}, SlotNotFound
		}
		return entities.SlotWithProviderEntity{}, err
	}

	return ent, nil
}

func (r *repository) DeleteSlot(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM availability_slots WHERE id = $1`, id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return SlotNotFound
	}

	return nil
}

func (r *repository) GetEncryptedUserKey(ctx context.Context, userID uuid.UUID) ([]byte, error) {
	var encryptedKey []byte
	err := r.db.QueryRowContext(ctx, `SELECT encrypted_user_key FROM users WHERE id = $1`, userID).Scan(&encryptedKey)
	if err != nil {
		return nil, err
	}
	return encryptedKey, nil
}

func (r *repository) AtomicBookSlot(ctx context.Context, id uuid.UUID, now time.Time) (int64, error) {
	result, err := r.db.ExecContext(ctx,
		`UPDATE availability_slots SET is_booked = true, updated_at = $2 WHERE id = $1 AND is_booked = false`,
		id, now,
	)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (r *repository) CreateAppointment(ctx context.Context, ent entities.AppointmentEntity) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO appointments (id, slot_id, veteran_id, title, status, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		ent.ID, ent.SlotID, ent.VeteranID, ent.Title, ent.Status, ent.CreatedAt, ent.UpdatedAt,
	)
	return err
}

func (r *repository) GetAppointmentWithSlot(ctx context.Context, id uuid.UUID) (entities.AppointmentWithSlotEntity, error) {
	var ent entities.AppointmentWithSlotEntity

	err := r.db.QueryRowContext(ctx,
		`SELECT a.id, a.slot_id, a.veteran_id, a.title, a.status, a.created_at, a.updated_at, s.provider_id, s.start_time, s.end_time
		 FROM appointments a
		 JOIN availability_slots s ON a.slot_id = s.id
		 WHERE a.id = $1`,
		id,
	).Scan(&ent.ID, &ent.SlotID, &ent.VeteranID, &ent.Title, &ent.Status, &ent.CreatedAt, &ent.UpdatedAt, &ent.SlotProviderID, &ent.StartTime, &ent.EndTime)

	if err == nil {
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.AppointmentWithSlotEntity{}, AppointmentNotFound
		}
		return entities.AppointmentWithSlotEntity{}, err
	}

	return ent, nil
}

func (r *repository) CancelAppointment(ctx context.Context, appointmentID uuid.UUID, now time.Time) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			slog.Error("failed to rollback transaction", "error", err)
		}
	}()

	var slotID uuid.UUID
	err = tx.QueryRowContext(ctx,
		`UPDATE appointments SET status = $2, updated_at = $3 WHERE id = $1 RETURNING slot_id`,
		appointmentID, string(entities.StatusCancelled), now,
	).Scan(&slotID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return AppointmentNotFound
		}
		return err
	}

	_, err = tx.ExecContext(ctx,
		`UPDATE availability_slots SET is_booked = false, updated_at = $2 WHERE id = $1`,
		slotID, now,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *repository) CheckSlotOverlap(ctx context.Context, providerID uuid.UUID, start, end time.Time) (bool, error) {
	var overlaps bool
	err := r.db.QueryRowContext(ctx,
		`SELECT EXISTS(
			SELECT 1 FROM availability_slots
			WHERE provider_id = $1 AND start_time < $3 AND end_time > $2
		)`,
		providerID, start, end,
	).Scan(&overlaps)
	return overlaps, err
}

func (r *repository) GetUrgentSlotMinutesForDate(ctx context.Context, providerID uuid.UUID, date time.Time) (int, error) {
	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
	dayEnd := dayStart.Add(24 * time.Hour)

	var minutes sql.NullFloat64
	err := r.db.QueryRowContext(ctx,
		`SELECT SUM(EXTRACT(EPOCH FROM (end_time - start_time)) / 60)
		 FROM availability_slots
		 WHERE provider_id = $1 AND is_urgent = true AND start_time >= $2 AND end_time <= $3`,
		providerID, dayStart, dayEnd,
	).Scan(&minutes)
	if err != nil {
		return 0, err
	}

	if !minutes.Valid {
		return 0, nil
	}

	return int(minutes.Float64), nil
}

func (r *repository) DeleteSlotsByProviderID(ctx context.Context, providerID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM availability_slots WHERE provider_id = $1`, providerID)
	return err
}

func (r *repository) GetAppointmentsByDay(ctx context.Context, veteranID uuid.UUID, dayStart, dayEnd time.Time) ([]entities.AppointmentWithSlotEntity, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT a.id, a.slot_id, a.veteran_id, a.title, a.status, a.created_at, a.updated_at, s.provider_id, s.start_time, s.end_time, u.encrypted_username, u.encrypted_first_name, u.encrypted_last_name, u.encrypted_user_key, u.image
		 FROM appointments a
		 JOIN availability_slots s ON a.slot_id = s.id
		 JOIN users u ON u.id = s.provider_id
		 WHERE a.veteran_id = $1
		   AND s.start_time >= $2
		   AND s.start_time < $3
		   AND a.status != 'CANCELLED'
		 ORDER BY s.start_time`,
		veteranID, dayStart, dayEnd,
	)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	var ents []entities.AppointmentWithSlotEntity
	for rows.Next() {
		var ent entities.AppointmentWithSlotEntity
		if err := rows.Scan(
			&ent.ID, &ent.SlotID, &ent.VeteranID, &ent.Title, &ent.Status, &ent.CreatedAt, &ent.UpdatedAt, &ent.SlotProviderID, &ent.StartTime, &ent.EndTime, &ent.ProviderUsernameEncrypted, &ent.ProviderFirstNameEncrypted, &ent.ProviderLastNameEncrypted, &ent.ProviderEncryptedUserKey, &ent.ProviderImage,
		); err != nil {
			return nil, err
		}
		ents = append(ents, ent)
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return ents, nil
}

func (r *repository) DeleteFutureSlotsBySeries(ctx context.Context, seriesID uuid.UUID, providerID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM availability_slots
		 WHERE series_id = $1 AND provider_id = $2 AND is_booked = false AND start_time > NOW()`,
		seriesID, providerID,
	)
	return err
}
