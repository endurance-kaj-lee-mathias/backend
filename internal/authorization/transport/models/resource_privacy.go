package models

type PatchResourcePrivacyRequest struct {
	IsPrivate bool `json:"isPrivate"`
}

func (r *PatchResourcePrivacyRequest) Validate() error {
	return nil
}

type ResourcePrivacyResponse struct {
	Resource  string `json:"resource"`
	IsPrivate bool   `json:"isPrivate"`
}
