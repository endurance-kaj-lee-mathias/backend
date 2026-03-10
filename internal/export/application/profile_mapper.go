package application

import (
	"strings"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/domain"
	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/export/infrastructure/entities"
)

func (s *service) decryptProfile(ent entities.UserExportEntity, userKey []byte) (domain.ProfileData, error) {
	email, err := s.enc.Decrypt(ent.EncryptedEmail, userKey)
	if err != nil {
		return domain.ProfileData{}, err
	}

	username, err := s.enc.Decrypt(ent.EncryptedUsername, userKey)
	if err != nil {
		return domain.ProfileData{}, err
	}

	firstName, err := s.enc.Decrypt(ent.EncryptedFirstName, userKey)
	if err != nil {
		return domain.ProfileData{}, err
	}

	lastName, err := s.enc.Decrypt(ent.EncryptedLastName, userKey)
	if err != nil {
		return domain.ProfileData{}, err
	}

	about, err := s.enc.Decrypt(ent.EncryptedAbout, userKey)
	if err != nil {
		return domain.ProfileData{}, err
	}

	introduction, err := s.enc.Decrypt(ent.EncryptedIntroduction, userKey)
	if err != nil {
		return domain.ProfileData{}, err
	}

	rolesStr, err := s.enc.Decrypt(ent.EncryptedRoles, userKey)
	if err != nil {
		return domain.ProfileData{}, err
	}
	roles := strings.Split(string(rolesStr), ",")
	if len(roles) == 1 && roles[0] == "" {
		roles = []string{}
	}

	var phoneNumber *string
	if len(ent.EncryptedPhoneNumber) > 0 {
		phone, err := s.enc.Decrypt(ent.EncryptedPhoneNumber, userKey)
		if err != nil {
			return domain.ProfileData{}, err
		}
		p := string(phone)
		phoneNumber = &p
	}

	return domain.ProfileData{
		Email:        string(email),
		Username:     string(username),
		FirstName:    string(firstName),
		LastName:     string(lastName),
		PhoneNumber:  phoneNumber,
		Roles:        roles,
		About:        string(about),
		Introduction: string(introduction),
		Image:        ent.Image,
		CreatedAt:    ent.CreatedAt,
		UpdatedAt:    ent.UpdatedAt,
	}, nil
}

func (s *service) decryptAddress(ent entities.AddressExportEntity, userKey []byte) (domain.AddressData, error) {
	street, err := s.enc.Decrypt(ent.EncryptedStreet, userKey)
	if err != nil {
		return domain.AddressData{}, err
	}

	locality, err := s.enc.Decrypt(ent.EncryptedLocality, userKey)
	if err != nil {
		return domain.AddressData{}, err
	}

	region, err := s.enc.Decrypt(ent.EncryptedRegion, userKey)
	if err != nil {
		return domain.AddressData{}, err
	}

	postalCode, err := s.enc.Decrypt(ent.EncryptedPostalCode, userKey)
	if err != nil {
		return domain.AddressData{}, err
	}

	country, err := s.enc.Decrypt(ent.EncryptedCountry, userKey)
	if err != nil {
		return domain.AddressData{}, err
	}

	return domain.AddressData{
		Street:     string(street),
		Locality:   string(locality),
		Region:     string(region),
		PostalCode: string(postalCode),
		Country:    string(country),
		CreatedAt:  ent.CreatedAt,
	}, nil
}
