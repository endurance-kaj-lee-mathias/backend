package infrastructure

import (
	"context"
	"log/slog"

	"github.com/gofrs/uuid"
)

func (r *veteranReader) GetVeteransForMember(ctx context.Context, memberID uuid.UUID) ([]VeteranProfile, error) {
	query := `
		SELECT u.id, u.encrypted_first_name, u.encrypted_last_name, u.encrypted_user_key, u.image
		FROM users u
		JOIN user_supports s ON u.id = s.veteran_id
		WHERE s.support_id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, memberID)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error("failed to close rows", "error", err)
		}
	}()

	var profiles []VeteranProfile

	for rows.Next() {
		var id uuid.UUID
		var encryptedFirst, encryptedLast, encryptedUserKey []byte
		var image *string

		if err := rows.Scan(&id, &encryptedFirst, &encryptedLast, &encryptedUserKey, &image); err != nil {
			return nil, err
		}

		userKey, err := r.enc.DecryptUserKey(encryptedUserKey)
		if err != nil {
			return nil, err
		}

		firstNameBytes, err := r.enc.Decrypt(encryptedFirst, userKey)
		if err != nil {
			return nil, err
		}

		lastNameBytes, err := r.enc.Decrypt(encryptedLast, userKey)
		if err != nil {
			return nil, err
		}

		imageStr := ""

		if image != nil {
			imageStr = *image
		}

		profiles = append(profiles, VeteranProfile{
			ID:        id,
			FirstName: string(firstNameBytes),
			LastName:  string(lastNameBytes),
			Image:     imageStr,
		})
	}

	return profiles, rows.Err()
}
