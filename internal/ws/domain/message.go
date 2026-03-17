package domain

import "time"

type OutboundMessage struct {
	Channel   string    `json:"channel"`
	SenderID  string    `json:"senderId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}
