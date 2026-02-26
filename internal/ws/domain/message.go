package domain

type MessageType string

const (
	MessageTypeSubscribe   MessageType = "subscribe"
	MessageTypeUnsubscribe MessageType = "unsubscribe"
	MessageTypeMessage     MessageType = "message"
)

type InboundMessage struct {
	Type    MessageType `json:"type"`
	Channel string      `json:"channel"`
	Payload any         `json:"payload"`
}

type OutboundMessage struct {
	Channel string `json:"channel"`
	From    string `json:"from"`
	Payload any    `json:"payload"`
}
