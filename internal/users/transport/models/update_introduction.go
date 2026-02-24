package models

type UpdateIntroductionModel struct {
	Introduction string `json:"introduction"`
}

func (m *UpdateIntroductionModel) Validate() error {
	if len(m.Introduction) > 500 {
		return InvalidIntroductionTooLong
	}
	return nil
}
