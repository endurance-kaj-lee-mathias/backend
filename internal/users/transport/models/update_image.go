package models

type UpdateImageModel struct {
	Image string `json:"image"`
}

func (m *UpdateImageModel) Validate() error {
	if m.Image == "" {
		return InvalidImageEmpty
	}
	return nil
}
