package models

type UpdatePhoneNumberModel struct {
	PhoneNumber *string `json:"phoneNumber"`
}

func (m *UpdatePhoneNumberModel) Validate() error {
	if *m.PhoneNumber == "" {
		return InvalidPhoneNumberEmpty
	}

	return nil
}
