package infrastructure

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/infrastructure/entities"
)

func (r *repository) GetUserWithAddress(ctx context.Context, userID uuid.UUID) (entities.UserExportEntity, *entities.AddressExportEntity, error) {
	query := `
		SELECT u.id, u.encrypted_email, u.encrypted_username, u.encrypted_first_name, u.encrypted_last_name,
			u.encrypted_phone_number, u.encrypted_roles, u.encrypted_about, u.encrypted_introduction,
			u.image, u.is_private, u.encrypted_user_key, u.created_at, u.updated_at,
			a.id, a.encrypted_street, a.encrypted_locality, a.encrypted_region,
			a.encrypted_postal_code, a.encrypted_country, a.created_at
		FROM users u
		LEFT JOIN addresses a ON a.user_id = u.id
		WHERE u.id = $1
	`

	var user entities.UserExportEntity
	var addrID *uuid.UUID
	var addrStreet, addrLocality, addrRegion, addrPostalCode, addrCountry []byte
	var addrCreatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.EncryptedEmail, &user.EncryptedUsername, &user.EncryptedFirstName, &user.EncryptedLastName,
		&user.EncryptedPhoneNumber, &user.EncryptedRoles, &user.EncryptedAbout, &user.EncryptedIntroduction,
		&user.Image, &user.IsPrivate, &user.EncryptedUserKey, &user.CreatedAt, &user.UpdatedAt,
		&addrID, &addrStreet, &addrLocality, &addrRegion,
		&addrPostalCode, &addrCountry, &addrCreatedAt,
	)
	if err != nil {
		return entities.UserExportEntity{}, nil, err
	}

	if addrID != nil {
		addr := &entities.AddressExportEntity{
			ID:                  *addrID,
			EncryptedStreet:     addrStreet,
			EncryptedLocality:   addrLocality,
			EncryptedRegion:     addrRegion,
			EncryptedPostalCode: addrPostalCode,
			EncryptedCountry:    addrCountry,
			CreatedAt:           addrCreatedAt.Time,
		}
		return user, addr, nil
	}

	return user, nil, nil
}

func (r *repository) GetStressSamples(ctx context.Context, userID uuid.UUID) ([]entities.StressSampleExportEntity, error) {
	query := `
		SELECT id, timestamp_utc, window_minutes, encrypted_mean_hr, encrypted_rmssd_ms,
			encrypted_resting_hr, encrypted_steps, encrypted_sleep_debt_hours, created_at
		FROM stress_samples
		WHERE user_id = $1
		ORDER BY timestamp_utc ASC
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

	var result []entities.StressSampleExportEntity
	for rows.Next() {
		var ent entities.StressSampleExportEntity
		if err := rows.Scan(
			&ent.ID, &ent.TimestampUTC, &ent.WindowMinutes,
			&ent.EncryptedMeanHR, &ent.EncryptedRMSSDms,
			&ent.EncryptedRestingHR, &ent.EncryptedSteps,
			&ent.EncryptedSleepDebtHours, &ent.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, ent)
	}

	return result, rows.Err()
}

func (r *repository) GetStressScores(ctx context.Context, userID uuid.UUID) ([]entities.StressScoreExportEntity, error) {
	query := `
		SELECT id, score, category, model_version, computed_at
		FROM stress_scores
		WHERE user_id = $1
		ORDER BY computed_at DESC
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

	var result []entities.StressScoreExportEntity
	for rows.Next() {
		var ent entities.StressScoreExportEntity
		if err := rows.Scan(&ent.ID, &ent.Score, &ent.Category, &ent.ModelVersion, &ent.ComputedAt); err != nil {
			return nil, err
		}
		result = append(result, ent)
	}

	return result, rows.Err()
}

func (r *repository) GetMoodEntries(ctx context.Context, userID uuid.UUID) ([]entities.MoodEntryExportEntity, error) {
	query := `
		SELECT id, date, mood_score, encrypted_notes, created_at, updated_at
		FROM mood_entries
		WHERE user_id = $1
		ORDER BY date DESC
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

	var result []entities.MoodEntryExportEntity
	for rows.Next() {
		var ent entities.MoodEntryExportEntity
		if err := rows.Scan(&ent.ID, &ent.Date, &ent.MoodScore, &ent.EncryptedNotes, &ent.CreatedAt, &ent.UpdatedAt); err != nil {
			return nil, err
		}
		result = append(result, ent)
	}

	return result, rows.Err()
}

func (r *repository) GetConversationsAndMessages(ctx context.Context, userID uuid.UUID) ([]entities.MessageExportEntity, error) {
	query := `
		SELECT m.id, m.conversation_id, m.sender_id, m.encrypted_content, m.created_at,
			cp_user.encrypted_conversation_key,
			cp_other.user_id,
			u_other.encrypted_username, u_other.encrypted_first_name, u_other.encrypted_last_name,
			u_other.encrypted_user_key
		FROM messages m
		JOIN conversation_participants cp_user ON cp_user.conversation_id = m.conversation_id AND cp_user.user_id = $1
		JOIN conversation_participants cp_other ON cp_other.conversation_id = m.conversation_id AND cp_other.user_id != $1
		JOIN users u_other ON u_other.id = cp_other.user_id
		WHERE m.conversation_id IN (
			SELECT conversation_id FROM conversation_participants WHERE user_id = $1
		)
		ORDER BY m.created_at ASC
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

	var result []entities.MessageExportEntity
	for rows.Next() {
		var ent entities.MessageExportEntity
		if err := rows.Scan(
			&ent.ID, &ent.ConversationID, &ent.SenderID, &ent.EncryptedContent, &ent.CreatedAt,
			&ent.EncryptedConversationKey,
			&ent.OtherParticipantID,
			&ent.OtherParticipantEncryptedUsername, &ent.OtherParticipantEncryptedFirstName,
			&ent.OtherParticipantEncryptedLastName, &ent.OtherParticipantEncryptedUserKey,
		); err != nil {
			return nil, err
		}
		result = append(result, ent)
	}

	return result, rows.Err()
}

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

func (r *repository) GetSupporters(ctx context.Context, veteranID uuid.UUID) ([]entities.SupportMemberExportEntity, error) {
	query := `
		SELECT u.id, u.encrypted_email, u.encrypted_username, u.encrypted_first_name, u.encrypted_last_name,
			u.image, u.encrypted_user_key, s.created_at
		FROM users u
		JOIN user_supports s ON u.id = s.support_id
		WHERE s.veteran_id = $1
	`

	return r.querySupportMembers(ctx, query, veteranID)
}

func (r *repository) GetSupportedVeterans(ctx context.Context, supportID uuid.UUID) ([]entities.SupportMemberExportEntity, error) {
	query := `
		SELECT u.id, u.encrypted_email, u.encrypted_username, u.encrypted_first_name, u.encrypted_last_name,
			u.image, u.encrypted_user_key, s.created_at
		FROM users u
		JOIN user_supports s ON u.id = s.veteran_id
		WHERE s.support_id = $1
	`

	return r.querySupportMembers(ctx, query, supportID)
}

func (r *repository) querySupportMembers(ctx context.Context, query string, id uuid.UUID) ([]entities.SupportMemberExportEntity, error) {
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	var result []entities.SupportMemberExportEntity
	for rows.Next() {
		var ent entities.SupportMemberExportEntity
		if err := rows.Scan(
			&ent.ID, &ent.EncryptedEmail, &ent.EncryptedUsername,
			&ent.EncryptedFirst, &ent.EncryptedLast,
			&ent.Image, &ent.EncryptedUserKey, &ent.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, ent)
	}

	return result, rows.Err()
}

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

func (r *repository) GetSentInvites(ctx context.Context, userID uuid.UUID) ([]entities.InviteExportEntity, error) {
	query := `
		SELECT i.id, i.receiver_id,
			r.encrypted_username, r.encrypted_first_name, r.encrypted_last_name, r.encrypted_user_key, r.image,
			i.status, i.note, i.created_at, i.updated_at
		FROM support_invites i
		JOIN users r ON r.id = i.receiver_id
		WHERE i.sender_id = $1
		ORDER BY i.created_at DESC
	`

	return r.queryInvites(ctx, query, userID)
}

func (r *repository) GetReceivedInvites(ctx context.Context, userID uuid.UUID) ([]entities.InviteExportEntity, error) {
	query := `
		SELECT i.id, i.sender_id,
			s.encrypted_username, s.encrypted_first_name, s.encrypted_last_name, s.encrypted_user_key, s.image,
			i.status, i.note, i.created_at, i.updated_at
		FROM support_invites i
		JOIN users s ON s.id = i.sender_id
		WHERE i.receiver_id = $1
		ORDER BY i.created_at DESC
	`

	return r.queryInvites(ctx, query, userID)
}

func (r *repository) queryInvites(ctx context.Context, query string, userID uuid.UUID) ([]entities.InviteExportEntity, error) {
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	var result []entities.InviteExportEntity
	for rows.Next() {
		var ent entities.InviteExportEntity
		if err := rows.Scan(
			&ent.ID, &ent.OtherUserID,
			&ent.OtherEncryptedUsername, &ent.OtherEncryptedFirstName, &ent.OtherEncryptedLastName,
			&ent.OtherEncryptedUserKey, &ent.OtherImage,
			&ent.Status, &ent.Note, &ent.CreatedAt, &ent.UpdatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, ent)
	}

	return result, rows.Err()
}
