package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/gofrs/uuid"
)

func (r *userRoleReader) GetRoles(ctx context.Context, userID uuid.UUID) ([]string, error) {
	var encryptedRoles []byte
	var encryptedUserKey []byte

	query := `SELECT encrypted_roles, encrypted_user_key FROM users WHERE id = $1`

	if err := r.db.QueryRowContext(ctx, query, userID).Scan(&encryptedRoles, &encryptedUserKey); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, UserNotFound
		}
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

	var roles []string
	if len(rolesBytes) == 0 {
		return roles, nil
	}

	if err := json.Unmarshal(rolesBytes, &roles); err != nil {
		return nil, err
	}

	return roles, nil
}
