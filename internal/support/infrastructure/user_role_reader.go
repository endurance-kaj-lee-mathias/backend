package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/gofrs/uuid"
)

func (r *userRoleReader) GetRole(ctx context.Context, userID uuid.UUID) (string, error) {
	var encryptedRoles []byte
	var encryptedUserKey []byte

	query := `SELECT encrypted_roles, encrypted_user_key FROM users WHERE id = $1`

	if err := r.db.QueryRowContext(ctx, query, userID).Scan(&encryptedRoles, &encryptedUserKey); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", UserNotFound
		}
		return "", err
	}

	userKey, err := r.enc.DecryptUserKey(encryptedUserKey)
	if err != nil {
		return "", err
	}

	rolesBytes, err := r.enc.Decrypt(encryptedRoles, userKey)
	if err != nil {
		return "", err
	}

	var roles []string
	if len(rolesBytes) == 0 {
		return "", nil
	}

	if err := json.Unmarshal(rolesBytes, &roles); err != nil {
		return "", err
	}

	if len(roles) > 0 {
		return roles[0], nil
	}

	return "", nil
}

func (r *userRoleReader) FindIDByUsername(ctx context.Context, username string) (uuid.UUID, error) {
	usernameHash := r.enc.Hash(username)

	var id uuid.UUID

	query := `SELECT id FROM users WHERE username_hash = $1`

	if err := r.db.QueryRowContext(ctx, query, usernameHash).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return uuid.UUID{}, UserNotFound
		}
		return uuid.UUID{}, err
	}

	return id, nil
}
