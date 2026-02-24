package models

type UpdateAboutModel struct {
	About string `json:"about"`
}

func (m *UpdateAboutModel) Validate() error {
	if len(m.About) > 500 {
		return InvalidAboutTooLong
	}
	return nil
}
