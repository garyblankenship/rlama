package agent

import (
	"fmt"
	"sync"
)

// SimpleMemory provides a basic in-memory implementation of the Memory interface
type SimpleMemory struct {
	store   map[string]interface{}
	history []string
	mu      sync.RWMutex
}

// NewSimpleMemory creates a new SimpleMemory instance
func NewSimpleMemory() *SimpleMemory {
	return &SimpleMemory{
		store:   make(map[string]interface{}),
		history: make([]string, 0),
	}
}

// Store stores a key-value pair in memory
func (m *SimpleMemory) Store(key string, value interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store[key] = value
	return nil
}

// Retrieve gets a value from memory by key
func (m *SimpleMemory) Retrieve(key string) (interface{}, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, exists := m.store[key]
	if !exists {
		return nil, fmt.Errorf("key %s not found in memory", key)
	}
	return value, nil
}

// GetHistory returns the conversation history
func (m *SimpleMemory) GetHistory() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Return a copy to prevent modification of internal slice
	history := make([]string, len(m.history))
	copy(history, m.history)
	return history
}

// AddToHistory adds an entry to conversation history
func (m *SimpleMemory) AddToHistory(entry string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.history = append(m.history, entry)
	return nil
}

// Clear clears all stored data and history
func (m *SimpleMemory) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.store = make(map[string]interface{})
	m.history = make([]string, 0)
}
