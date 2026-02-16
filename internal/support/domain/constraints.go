package domain

const (
	RoleVeteran   = "VETERAN"
	RoleTherapist = "THERAPIST"
	RoleSupport   = "SUPPORT"
)

func ValidateSupportRelationship(veteranRoles, supporterRoles []string, veteranID, supporterID string) error {
	if veteranID == supporterID {
		return ErrSelfSupportNotAllowed
	}

	if !hasRole(veteranRoles, RoleVeteran) {
		return ErrVeteranMustHaveVeteranRole
	}

	if !canSupport(supporterRoles) {
		return ErrSupporterMustBeAbleToSupport
	}

	return nil
}

func hasRole(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func canSupport(roles []string) bool {
	return hasRole(roles, RoleVeteran) || hasRole(roles, RoleTherapist) || hasRole(roles, RoleSupport)
}
