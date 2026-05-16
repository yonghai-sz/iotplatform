package mqtt

import (
	"sync"

	gomqtt "github.com/eclipse/paho.mqtt.golang"
)

// Manager holds a connected MQTT client and its subscription list.
type Manager struct {
	mu sync.Mutex

	client            gomqtt.Client
	disconnectQuiesce uint
}

// NewManager dials the broker once
// and registers onConnect resubscribe for all items.
func NewManager(dial DialOptions, onConnect func(gomqtt.Client), disconnectQuiesce uint) (*Manager, error) {

	m := &Manager{
		disconnectQuiesce: disconnectQuiesce,
	}

	client, err := connectClient(dial, onConnect)
	if err != nil {
		return nil, err
	}

	m.client = client
	return m, nil
}

// Client returns the underlying paho client, or nil if not connected.
func (m *Manager) Client() gomqtt.Client {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.client
}

// Disconnect cleanly closes the MQTT session.
func (m *Manager) Disconnect() {
	m.mu.Lock()
	c := m.client
	ms := m.disconnectQuiesce
	m.mu.Unlock()

	if c == nil {
		return
	}
	if ms == 0 {
		ms = 250
	}
	c.Disconnect(ms)
}
