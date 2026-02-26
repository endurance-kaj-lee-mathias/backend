package application

import "gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/ws/domain"

const sendBufferSize = 16

type Client struct {
	UserID string
	send   chan domain.OutboundMessage
}

func newClient(userID string) *Client {
	return &Client{
		UserID: userID,
		send:   make(chan domain.OutboundMessage, sendBufferSize),
	}
}

func (c *Client) Send(msg domain.OutboundMessage) {
	select {
	case c.send <- msg:
	default:
	}
}

func (c *Client) Receive() <-chan domain.OutboundMessage {
	return c.send
}

func (c *Client) Close() {
	close(c.send)
}
