package infrastructure

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/gofrs/uuid"
)

func (r *userRoleReader) GetRoles(ctx context.Context, userID uuid.UUID) ([]string, error) {
	var rawRoles []byte

	query := `SELECT roles FROM users WHERE id = $1`

	if err := r.db.QueryRowContext(ctx, query, userID).Scan(&rawRoles); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	var roles []string
	if len(rawRoles) == 0 {
		return roles, nil
	}

	if err := json.Unmarshal(rawRoles, &roles); err != nil {
		return nil, err
	}

	return roles, nil
}
