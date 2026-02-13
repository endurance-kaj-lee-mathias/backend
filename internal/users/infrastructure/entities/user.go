package entities

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
)

type UserEntity struct {
	ID        uuid.UUID       `db:"id"`
	Email     string          `db:"email"`
	Roles     json.RawMessage `db:"roles"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
}

func FromEntity(entity UserEntity) (domain.User, error) {
	roles := make([]domain.Role, 0)

	if len(entity.Roles) > 0 {
		if err := json.Unmarshal(entity.Roles, &roles); err != nil {
			return domain.User{}, fmt.Errorf("unmarshal roles: %w", err)
		}
	}

	return domain.User{
		ID:        entity.ID,
		Email:     entity.Email,
		Roles:     roles,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
	}, nil
}
