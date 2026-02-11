package models

import "gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/message/domain"

type MessageModel struct {
	Value string `json:"value"`
}

func ToModel(msg domain.Message) MessageModel {
	return MessageModel{
		Value: msg.Value,
	}
}
