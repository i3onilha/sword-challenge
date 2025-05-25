package messaging

import (
	"context"
	"sync"
)

// MockBroker is a simple in-memory message broker for testing
type MockBroker struct {
	mu       sync.RWMutex
	messages []TaskCreatedMessage
}

type TaskCreatedMessage struct {
	TaskID       int64
	TechnicianID int64
	Title        string
}

// NewMockBroker creates a new mock message broker
func NewMockBroker() *MockBroker {
	return &MockBroker{
		messages: make([]TaskCreatedMessage, 0),
	}
}

// PublishTaskCreated implements MessageBroker interface
func (m *MockBroker) PublishTaskCreated(ctx context.Context, taskID int64, technicianID int64, title string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.messages = append(m.messages, TaskCreatedMessage{
		TaskID:       taskID,
		TechnicianID: technicianID,
		Title:        title,
	})
	return nil
}

// Close implements MessageBroker interface
func (m *MockBroker) Close() error {
	return nil
}

// GetMessages returns all published messages
func (m *MockBroker) GetMessages() []TaskCreatedMessage {
	m.mu.RLock()
	defer m.mu.RUnlock()

	messages := make([]TaskCreatedMessage, len(m.messages))
	copy(messages, m.messages)
	return messages
}

// ClearMessages clears all published messages
func (m *MockBroker) ClearMessages() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.messages = make([]TaskCreatedMessage, 0)
}
