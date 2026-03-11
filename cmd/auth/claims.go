package auth

type ClaimsAddress struct {
	StreetAddress string `json:"street_address"`
	Locality      string `json:"locality"`
	Region        string `json:"region"`
	PostalCode    string `json:"postal_code"`
	Country       string `json:"country"`
}

type Claims struct {
	Sub         string        `json:"sub"`
	Email       string        `json:"email"`
	Username    string        `json:"preferred_username"`
	FirstName   string        `json:"given_name"`
	LastName    string        `json:"family_name"`
	PhoneNumber string        `json:"phoneNumber"`
	Address     ClaimsAddress `json:"address"`
	Roles       []string      `json:"roles"`
	ClientID    string        `json:"azp"`
}
