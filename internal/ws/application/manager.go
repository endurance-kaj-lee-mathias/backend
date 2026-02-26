package application

import (
	"sync"

	"gitlab.com/kdg-ti/the-lab/teams-25-26/26-de-uitgeruste-it-ers/backend/internal/ws/domain"
)

type Manager struct {
	mu       sync.RWMutex
	channels map[string]map[*Client]struct{}
}

func NewManager() *Manager {
	return &Manager{
		channels: make(map[string]map[*Client]struct{}),
	}
}

func (m *Manager) Subscribe(channelID string, client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.channels[channelID] == nil {
		m.channels[channelID] = make(map[*Client]struct{})
	}

	m.channels[channelID][client] = struct{}{}
}

func (m *Manager) Unsubscribe(channelID string, client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	clients, ok := m.channels[channelID]
	if !ok {
		return
	}

	delete(clients, client)

	if len(clients) == 0 {
		delete(m.channels, channelID)
	}
}

func (m *Manager) UnsubscribeAll(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for channelID, clients := range m.channels {
		delete(clients, client)

		if len(clients) == 0 {
			delete(m.channels, channelID)
		}
	}
}

func (m *Manager) Broadcast(channelID string, msg domain.OutboundMessage) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for client := range m.channels[channelID] {
		client.Send(msg)
	}
}

func NewClient(userID string) *Client {
	return newClient(userID)
}

func (m *Manager) GetChannelSubscribers(channelID string) int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.channels[channelID])
}

func (m *Manager) GetChannels() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.channels)
}
