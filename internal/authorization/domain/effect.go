package domain

type PolicyEffect string

const (
	EffectAllow PolicyEffect = "allow"
	EffectDeny  PolicyEffect = "deny"
)

func ValidEffect(e string) bool {
	switch PolicyEffect(e) {
	case EffectAllow, EffectDeny:
		return true
	}
	return false
}
