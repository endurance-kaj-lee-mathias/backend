package domain

type ResourceType string

const (
	ResourceUserProfile  ResourceType = "userProfile"
	ResourceStressScores ResourceType = "stressScores"
	ResourceMoodEntries  ResourceType = "moodEntries"
	ResourceCalendar     ResourceType = "calendar"
)

func ValidResource(r string) bool {
	switch ResourceType(r) {
	case ResourceUserProfile, ResourceStressScores, ResourceMoodEntries, ResourceCalendar:
		return true
	}
	return r == "*"
}
