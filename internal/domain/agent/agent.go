package agent

import (
	"context"
)

// AgentMode defines the operation mode of an agent
type AgentMode string

const (
	// ConversationalMode represents an agent that operates in a conversational manner
	ConversationalMode AgentMode = "conversational"
	// AutonomousMode represents an agent that operates autonomously
	AutonomousMode AgentMode = "autonomous"
	// OrchestratedMode represents an agent that uses task orchestration for complex queries
	OrchestratedMode AgentMode = "orchestrated"
)

// Tool interface is defined in tool.go

// Memory represents the agent's memory storage
type Memory interface {
	// Store stores a key-value pair in memory
	Store(key string, value interface{}) error
	// Retrieve gets a value from memory by key
	Retrieve(key string) (interface{}, error)
	// GetHistory returns the conversation history
	GetHistory() []string
	// AddToHistory adds an entry to conversation history
	AddToHistory(entry string) error
}

// Agent represents an intelligent agent that can use tools to accomplish tasks
type Agent interface {
	// Run executes the agent with a given query/goal
	Run(ctx context.Context, input string) (string, error)
	// AddTool adds a new tool to the agent's capabilities
	AddTool(tool Tool) error
	// GetTools returns all tools available to the agent
	GetTools() []Tool
	// GetMode returns the agent's operation mode
	GetMode() AgentMode
	// GetMemory returns the agent's memory
	GetMemory() Memory
}
