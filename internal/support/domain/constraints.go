package domain

func ValidateSupportRelationship(_ string, _ string, veteranID, supporterID string) error {
	if veteranID == supporterID {
		return SelfSupportNotAllowed
	}

	return nil
}
