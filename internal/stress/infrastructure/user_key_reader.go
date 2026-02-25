package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gofrs/uuid"
)

func (r *userKeyReader) GetEncryptedUserKey(ctx context.Context, userID uuid.UUID) ([]byte, error) {
	var encryptedUserKey []byte

	query := `SELECT encrypted_user_key FROM users WHERE id = $1`

	if err := r.db.QueryRowContext(ctx, query, userID).Scan(&encryptedUserKey); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, UserNotFound
		}
		return nil, err
	}

	return encryptedUserKey, nil
}
