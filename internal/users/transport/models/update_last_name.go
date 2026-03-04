package models

type UpdateLastNameModel struct {
	LastName string `json:"lastName"`
}

func (m *UpdateLastNameModel) Validate() error {
	if m.LastName == "" {
		return InvalidLastNameEmpty
	}
	return nil
}
