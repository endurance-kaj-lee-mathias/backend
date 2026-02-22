package entities

import (
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
)

type UserEntity struct {
	ID          uuid.UUID       `db:"id"`
	Email       string          `db:"email"`
	FirstName   string          `db:"first_name"`
	LastName    string          `db:"last_name"`
	PhoneNumber *string         `db:"phone_number"`
	Roles       json.RawMessage `db:"roles"`
	CreatedAt   time.Time       `db:"created_at"`
	UpdatedAt   time.Time       `db:"updated_at"`
}

func FromEntity(ent UserEntity) (domain.User, error) {
	roles := make([]domain.Role, 0)

	if len(ent.Roles) > 0 {
		if err := json.Unmarshal(ent.Roles, &roles); err != nil {
			return domain.User{}, InvalidRoles
		}
	}

	id, err := domain.ParseId(ent.ID.String())

	if err != nil {
		return domain.User{}, err
	}

	return domain.User{
		ID:          id,
		Email:       ent.Email,
		FirstName:   ent.FirstName,
		LastName:    ent.LastName,
		PhoneNumber: ent.PhoneNumber,
		Roles:       roles,
		CreatedAt:   ent.CreatedAt,
		UpdatedAt:   ent.UpdatedAt,
	}, nil
}

func FromEntities(entities []UserEntity) []domain.User {
	output := make([]domain.User, 0, len(entities))

	for _, entity := range entities {
		entity, err := FromEntity(entity)

		if err != nil {
			continue
		}

		output = append(output, entity)
	}

	return output
}

func ToEntity(usr domain.User) (UserEntity, error) {
	roles, err := json.Marshal(usr.Roles)

	if err != nil {
		return UserEntity{}, InvalidRoles
	}

	return UserEntity{
		ID:          usr.ID.UUID,
		Email:       usr.Email,
		FirstName:   usr.FirstName,
		LastName:    usr.LastName,
		PhoneNumber: usr.PhoneNumber,
		Roles:       roles,
		CreatedAt:   usr.CreatedAt,
		UpdatedAt:   usr.UpdatedAt,
	}, nil
}
