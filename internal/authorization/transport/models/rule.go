package models

import (
	"time"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/authorization/domain"
)

type CreateRuleRequest struct {
	ViewerID string `json:"viewerId"`
	Resource string `json:"resource"`
	Effect   string `json:"effect"`
}

func (r *CreateRuleRequest) Validate() error {
	if r.ViewerID == "" {
		return InvalidViewerID
	}

	if !domain.ValidResource(r.Resource) {
		return InvalidResource
	}

	if !domain.ValidEffect(r.Effect) {
		return InvalidEffect
	}

	return nil
}

type RuleResponse struct {
	ID        uuid.UUID `json:"id"`
	OwnerID   uuid.UUID `json:"ownerId"`
	ViewerID  uuid.UUID `json:"viewerId"`
	Resource  string    `json:"resource"`
	Effect    string    `json:"effect"`
	CreatedAt time.Time `json:"createdAt"`
}

func ToRuleResponse(r domain.Rule) RuleResponse {
	return RuleResponse{
		ID:        r.ID,
		OwnerID:   r.OwnerID,
		ViewerID:  r.ViewerID,
		Resource:  string(r.Resource),
		Effect:    string(r.Effect),
		CreatedAt: r.CreatedAt,
	}
}

func ToRuleResponses(rules []domain.Rule) []RuleResponse {
	out := make([]RuleResponse, len(rules))
	for i, r := range rules {
		out[i] = ToRuleResponse(r)
	}
	return out
}
