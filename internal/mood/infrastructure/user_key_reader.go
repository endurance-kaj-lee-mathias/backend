package infrastructure

import (
	"context"
	"database/sql"
	"errors"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/mood/domain"
)

func (r *userKeyReader) GetEncryptedUserKey(ctx context.Context, userID domain.UserId) ([]byte, error) {
	var encryptedUserKey []byte

	query := `SELECT encrypted_user_key FROM users WHERE id = $1`

	if err := r.db.QueryRowContext(ctx, query, userID.UUID).Scan(&encryptedUserKey); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, UserNotFound
		}
		return nil, err
	}

	return encryptedUserKey, nil
}
