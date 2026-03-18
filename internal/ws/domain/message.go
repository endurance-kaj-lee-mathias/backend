package domain

import "time"

type OutboundMessage struct {
	Channel   string    `json:"channel"`
	SenderID  string    `json:"senderId"`
	Username  string    `json:"username"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
}
