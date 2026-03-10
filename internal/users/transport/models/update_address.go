package models

type UpdateAddressModel struct {
	Street     string `json:"street"`
	Locality   string `json:"locality"`
	Region     string `json:"region"`
	PostalCode string `json:"postalCode"`
	Country    string `json:"country"`
}

func (m *UpdateAddressModel) Validate() error {
	if m.Street == "" {
		return InvalidStreet
	}

	if m.Locality == "" {
		return InvalidLocality
	}

	if m.PostalCode == "" {
		return InvalidPostalCode
	}

	if m.Country == "" {
		return InvalidCountry
	}

	return nil
}
