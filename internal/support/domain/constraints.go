package domain

import usersdomain "gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"

func ValidateSupportRelationship(veteranRoles, supporterRoles []string, veteranID, supporterID string) error {
	if veteranID == supporterID {
		return ErrSelfSupportNotAllowed
	}

	if !hasRole(veteranRoles, string(usersdomain.RoleVeteran)) {
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
	return hasRole(roles, string(usersdomain.RoleVeteran)) ||
		hasRole(roles, string(usersdomain.RoleTherapist)) ||
		hasRole(roles, string(usersdomain.RoleSupport))
}
