package mqtt

import (
	"fmt"
	"time"
)

// Publish publishes a message to the given topic.
func (m *Manager) Publish(topic string, qos byte, retained bool, payload any) error {

	m.mu.Lock()
	client := m.client
	m.mu.Unlock()

	if client == nil {
		return fmt.Errorf("mqtt client is nil")
	}

	token := client.Publish(topic, qos, retained, payload)
	if token.Error() != nil {
		return fmt.Errorf("pub message to topic %s error: %w", topic, token.Error())
	}

	ack := token.WaitTimeout(3 * time.Second)
	if !ack {
		return fmt.Errorf("pub message to topic %s timeout", topic)
	}

	return nil
}
