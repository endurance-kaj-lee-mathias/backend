package models

import (
	"errors"

	"github.com/gofrs/uuid"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/users/domain"
)

type UserModel struct {
	ID               uuid.UUID              `json:"id"`
	FirstName        string                 `json:"firstName"`
	LastName         string                 `json:"lastName"`
	Username         string                 `json:"username"`
	PhoneNumber      *string                `json:"phoneNumber,omitempty"`
	About            string                 `json:"about"`
	Introduction     string                 `json:"introduction"`
	Image            string                 `json:"image"`
	RiskLevel        string                 `json:"riskLevel"`
	IsPrivate        bool                   `json:"isPrivate"`
	Address          *AddressModel          `json:"address,omitempty"`
	SharingResources []ResourcePrivacyModel `json:"sharingResources"`
}

type ResourcePrivacyModel struct {
	Resource  string `json:"resource"`
	IsPrivate bool   `json:"isPrivate"`
}

func ToModel(usr domain.User, addr *domain.Address, sharingResources map[string]bool) UserModel {
	m := UserModel{
		ID:           usr.ID.UUID,
		FirstName:    usr.FirstName,
		LastName:     usr.LastName,
		Username:     usr.Username,
		PhoneNumber:  usr.PhoneNumber,
		About:        usr.About,
		Introduction: usr.Introduction,
		Image:        usr.Image,
		RiskLevel:    string(usr.RiskLevel),
		IsPrivate:    usr.IsPrivate,
	}
	if addr != nil {
		a := ToAddressModel(*addr)
		m.Address = &a
	}

	if sharingResources != nil {
		var resources []ResourcePrivacyModel
		for res, isPrivate := range sharingResources {
			resources = append(resources, ResourcePrivacyModel{
				Resource:  res,
				IsPrivate: isPrivate,
			})
		}
		if resources != nil {
			m.SharingResources = resources
		} else {
			m.SharingResources = []ResourcePrivacyModel{}
		}
	} else {
		m.SharingResources = []ResourcePrivacyModel{}
	}

	return m
}

type UpdateRiskLevelModel struct {
	RiskLevel string `json:"riskLevel"`
}

func (m UpdateRiskLevelModel) Validate() error {
	if m.RiskLevel != string(domain.RiskLevelNormal) && m.RiskLevel != string(domain.RiskLevelHigh) {
		return errors.New("invalid risk level")
	}
	return nil
}
