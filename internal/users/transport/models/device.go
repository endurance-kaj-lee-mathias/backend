package models

import "errors"

type UpsertDeviceModel struct {
	Token    string `json:"token"`
	Platform string `json:"platform"`
}

func (m *UpsertDeviceModel) Validate() error {
	if m.Token == "" {
		return errors.New("token is required")
	}
	if m.Platform != "ios" && m.Platform != "android" {
		return errors.New("platform must be ios or android")
	}
	return nil
}

type DeleteDeviceModel struct {
	Token string `json:"token"`
}

func (m *DeleteDeviceModel) Validate() error {
	if m.Token == "" {
		return errors.New("token is required")
	}
	return nil
}
