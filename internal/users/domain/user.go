package domain

import (
	"strings"
	"time"
)

type User struct {
	ID        UserId    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Roles     []Role    `json:"roles"`
	Phone     *string   `json:"phone,omitempty"`
	Address   *Address  `json:"address,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func NewUser(id UserId, email string, firstName string, lastName string, roles []Role) User {
	now := time.Now().UTC()

	return User{
		ID:        id,
		Email:     email,
		FirstName: firstName,
		LastName:  lastName,
		Roles:     roles,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (u *User) UpdateProfile(phone *string, address *Address) error {
	if phone != nil {
		normalized, err := normalizePhone(*phone)
		if err != nil {
			return err
		}
		u.Phone = &normalized
	}

	if address != nil {
		if err := address.Validate(); err != nil {
			return err
		}
		addr := *address
		u.Address = &addr
	}

	if phone != nil || address != nil {
		u.UpdatedAt = time.Now().UTC()
	}

	return nil
}

func normalizePhone(value string) (string, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", InvalidPhone
	}

	var b strings.Builder
	for i, r := range trimmed {
		switch {
		case r == '+' && i == 0:
			b.WriteRune(r)
		case r >= '0' && r <= '9':
			b.WriteRune(r)
		case r == ' ' || r == '-' || r == '(' || r == ')' || r == '.':
			continue
		default:
			return "", InvalidPhone
		}
	}

	normalized := b.String()
	digits := 0
	for _, r := range normalized {
		if r >= '0' && r <= '9' {
			digits++
		}
	}

	if digits < 9 || digits > 15 {
		return "", InvalidPhone
	}

	return normalized, nil
}
