package auth

type Claims struct {
	Sub       string   `json:"sub"`
	Email     string   `json:"email"`
	FirstName string   `json:"given_name"`
	LastName  string   `json:"family_name"`
	Roles     []string `json:"roles"`
}
