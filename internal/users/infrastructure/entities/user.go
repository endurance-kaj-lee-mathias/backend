package entities

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
)

type UserEntity struct {
	ID        uuid.UUID       `db:"id"`
	Email     string          `db:"email"`
	FirstName string          `db:"first_name"`
	LastName  string          `db:"last_name"`
	Roles     json.RawMessage `db:"roles"`
	Phone     sql.NullString  `db:"phone"`
	Street    sql.NullString  `db:"street"`
	Number    sql.NullString  `db:"number"`
	Postal    sql.NullString  `db:"postal"`
	City      sql.NullString  `db:"city"`
	Country   sql.NullString  `db:"country"`
	CreatedAt time.Time       `db:"created_at"`
	UpdatedAt time.Time       `db:"updated_at"`
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

	var phone *string
	if ent.Phone.Valid {
		val := ent.Phone.String
		phone = &val
	}

	var address *domain.Address
	if ent.Street.Valid && ent.Number.Valid && ent.Postal.Valid && ent.City.Valid && ent.Country.Valid {
		addr := domain.Address{
			Street:  ent.Street.String,
			Number:  ent.Number.String,
			Postal:  ent.Postal.String,
			City:    ent.City.String,
			Country: ent.Country.String,
		}
		address = &addr
	}

	return domain.User{
		ID:        id,
		Email:     ent.Email,
		FirstName: ent.FirstName,
		LastName:  ent.LastName,
		Roles:     roles,
		Phone:     phone,
		Address:   address,
		CreatedAt: ent.CreatedAt,
		UpdatedAt: ent.UpdatedAt,
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

	ent := UserEntity{
		ID:        usr.ID.UUID,
		Email:     usr.Email,
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
		Roles:     roles,
		CreatedAt: usr.CreatedAt,
		UpdatedAt: usr.UpdatedAt,
	}

	if usr.Phone != nil {
		ent.Phone = sql.NullString{String: *usr.Phone, Valid: true}
	}

	if usr.Address != nil {
		ent.Street = sql.NullString{String: usr.Address.Street, Valid: true}
		ent.Number = sql.NullString{String: usr.Address.Number, Valid: true}
		ent.Postal = sql.NullString{String: usr.Address.Postal, Valid: true}
		ent.City = sql.NullString{String: usr.Address.City, Valid: true}
		ent.Country = sql.NullString{String: usr.Address.Country, Valid: true}
	}

	return ent, nil
}
