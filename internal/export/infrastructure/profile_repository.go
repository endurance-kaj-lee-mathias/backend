package infrastructure

import (
	"context"
	"database/sql"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/infrastructure/entities"
)

func (r *repository) GetUserWithAddress(ctx context.Context, userID uuid.UUID) (entities.UserExportEntity, *entities.AddressExportEntity, error) {
	query := `
		SELECT u.id, u.encrypted_email, u.encrypted_username, u.encrypted_first_name, u.encrypted_last_name,
			u.encrypted_phone_number, u.encrypted_roles, u.encrypted_about, u.encrypted_introduction,
			u.image, u.is_private, u.encrypted_user_key, u.created_at, u.updated_at,
			a.id, a.encrypted_street, a.encrypted_locality, a.encrypted_region,
			a.encrypted_postal_code, a.encrypted_country, a.created_at
		FROM users u
		LEFT JOIN user_addresses a ON a.user_id = u.id
		WHERE u.id = $1
	`

	var user entities.UserExportEntity
	var addrID *uuid.UUID
	var addrStreet, addrLocality, addrRegion, addrPostalCode, addrCountry []byte
	var addrCreatedAt sql.NullTime

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.EncryptedEmail, &user.EncryptedUsername, &user.EncryptedFirstName, &user.EncryptedLastName,
		&user.EncryptedPhoneNumber, &user.EncryptedRoles, &user.EncryptedAbout, &user.EncryptedIntroduction,
		&user.Image, &user.IsPrivate, &user.EncryptedUserKey, &user.CreatedAt, &user.UpdatedAt,
		&addrID, &addrStreet, &addrLocality, &addrRegion,
		&addrPostalCode, &addrCountry, &addrCreatedAt,
	)
	if err != nil {
		return entities.UserExportEntity{}, nil, err
	}

	if addrID != nil {
		addr := &entities.AddressExportEntity{
			ID:                  *addrID,
			EncryptedStreet:     addrStreet,
			EncryptedLocality:   addrLocality,
			EncryptedRegion:     addrRegion,
			EncryptedPostalCode: addrPostalCode,
			EncryptedCountry:    addrCountry,
			CreatedAt:           addrCreatedAt.Time,
		}
		return user, addr, nil
	}

	return user, nil, nil
}
