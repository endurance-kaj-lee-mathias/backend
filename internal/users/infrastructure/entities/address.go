package entities

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
)

type AddressEntity struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	Street      string    `db:"street"`
	HouseNumber string    `db:"house_number"`
	PostalCode  string    `db:"postal_code"`
	City        string    `db:"city"`
	Country     string    `db:"country"`
	CreatedAt   time.Time `db:"created_at"`
}

func AddressFromEntity(ent AddressEntity) (domain.Address, error) {
	id, err := domain.ParseAddressId(ent.ID.String())
	if err != nil {
		return domain.Address{}, err
	}

	userID, err := domain.ParseId(ent.UserID.String())
	if err != nil {
		return domain.Address{}, err
	}

	return domain.Address{
		ID:          id,
		UserID:      userID,
		Street:      ent.Street,
		HouseNumber: ent.HouseNumber,
		PostalCode:  ent.PostalCode,
		City:        ent.City,
		Country:     ent.Country,
		CreatedAt:   ent.CreatedAt,
	}, nil
}

func AddressToEntity(a domain.Address) AddressEntity {
	return AddressEntity{
		ID:          a.ID.UUID,
		UserID:      a.UserID.UUID,
		Street:      a.Street,
		HouseNumber: a.HouseNumber,
		PostalCode:  a.PostalCode,
		City:        a.City,
		Country:     a.Country,
		CreatedAt:   a.CreatedAt,
	}
}
