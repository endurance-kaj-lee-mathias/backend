package entities

import (
	"time"

	"github.com/gofrs/uuid"
)

type UserDeviceEntity struct {
	ID          uuid.UUID `db:"id"`
	UserID      uuid.UUID `db:"user_id"`
	DeviceToken string    `db:"device_token"`
	Platform    string    `db:"platform"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
