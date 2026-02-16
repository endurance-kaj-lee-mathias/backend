package domain

import "gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"

func ValidateSupportRelationship(veteranRoles, supporterRoles []string, veteranID, supporterID string) error {
	if veteranID == supporterID {
		return SelfSupportNotAllowed
	}

	if !hasRole(veteranRoles, string(domain.RoleVeteran)) {
		return VeteranMustHaveVeteranRole
	}

	if !canSupport(supporterRoles) {
		return SupporterMustBeAbleToSupport
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
	return hasRole(roles, string(domain.RoleVeteran)) ||
		hasRole(roles, string(domain.RoleTherapist)) ||
		hasRole(roles, string(domain.RoleSupport))
}
