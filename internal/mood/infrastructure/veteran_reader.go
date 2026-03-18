package infrastructure

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/gofrs/uuid"
)

func (r *veteranReader) GetVeteransForMember(ctx context.Context, memberID uuid.UUID) ([]VeteranProfile, error) {
	query := `
		SELECT u.id, u.encrypted_username, u.encrypted_first_name, u.encrypted_last_name,
			u.encrypted_user_key, u.image, u.encrypted_roles
		FROM users u
		JOIN user_supports s ON u.id = s.veteran_id
		WHERE s.support_id = $1
		UNION
		SELECT u.id, u.encrypted_username, u.encrypted_first_name, u.encrypted_last_name,
			u.encrypted_user_key, u.image, u.encrypted_roles
		FROM users u
		JOIN user_supports s ON u.id = s.support_id
		WHERE s.veteran_id = $1
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
		var encryptedUsername, encryptedFirst, encryptedLast, encryptedUserKey, encryptedRoles []byte
		var image *string

		if err := rows.Scan(&id, &encryptedUsername, &encryptedFirst, &encryptedLast, &encryptedUserKey, &image, &encryptedRoles); err != nil {
			return nil, err
		}

		userKey, err := r.enc.DecryptUserKey(encryptedUserKey)
		if err != nil {
			return nil, err
		}

		rolesBytes, err := r.enc.Decrypt(encryptedRoles, userKey)
		if err != nil {
			return nil, err
		}

		if !hasVeteranRole(rolesBytes) {
			continue
		}

		usernameBytes, err := r.enc.Decrypt(encryptedUsername, userKey)
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
			Username:  string(usernameBytes),
			FirstName: string(firstNameBytes),
			LastName:  string(lastNameBytes),
			Image:     imageStr,
		})
	}

	return profiles, rows.Err()
}

func hasVeteranRole(rolesBytes []byte) bool {
	if len(rolesBytes) == 0 {
		return false
	}

	var roles []string
	if err := json.Unmarshal(rolesBytes, &roles); err != nil {
		return false
	}

	for _, role := range roles {
		if role == "veteran" {
			return true
		}
	}

	return false
}
