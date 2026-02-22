package domain

import "strings"

type Address struct {
	Street  string `json:"street"`
	Number  string `json:"number"`
	Postal  string `json:"postal"`
	City    string `json:"city"`
	Country string `json:"country"`
}

func NewAddress(street string, number string, postal string, city string, country string) (Address, error) {
	addr := Address{
		Street:  street,
		Number:  number,
		Postal:  postal,
		City:    city,
		Country: country,
	}

	if err := addr.Validate(); err != nil {
		return Address{}, err
	}

	return addr, nil
}

func (a Address) Validate() error {
	if strings.TrimSpace(a.Street) == "" ||
		strings.TrimSpace(a.Number) == "" ||
		strings.TrimSpace(a.Postal) == "" ||
		strings.TrimSpace(a.City) == "" ||
		strings.TrimSpace(a.Country) == "" {
		return InvalidAddress
	}

	return nil
}
