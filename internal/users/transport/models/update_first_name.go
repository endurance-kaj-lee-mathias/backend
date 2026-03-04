package models

type UpdateFirstNameModel struct {
	FirstName string `json:"firstName"`
}

func (m *UpdateFirstNameModel) Validate() error {
	if m.FirstName == "" {
		return InvalidFirstNameEmpty
	}
	return nil
}
