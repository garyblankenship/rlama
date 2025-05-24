package agent

import (
	"context"
	"fmt"
	"sync"
)

// baseAgent provides a base implementation of the Agent interface
type baseAgent struct {
	mode   AgentMode
	tools  []Tool
	memory Memory
	mu     sync.RWMutex
}

// NewBaseAgent creates a new base agent with the given mode
func NewBaseAgent(mode AgentMode, memory Memory) *baseAgent {
	return &baseAgent{
		mode:   mode,
		tools:  make([]Tool, 0),
		memory: memory,
	}
}

// AddTool adds a new tool to the agent's capabilities
func (a *baseAgent) AddTool(tool Tool) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Check if tool with same name already exists
	for _, t := range a.tools {
		if t.Name() == tool.Name() {
			return fmt.Errorf("tool with name %s already exists", tool.Name())
		}
	}

	a.tools = append(a.tools, tool)
	return nil
}

// GetTools returns all tools available to the agent
func (a *baseAgent) GetTools() []Tool {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// Return a copy to prevent modification of internal slice
	tools := make([]Tool, len(a.tools))
	copy(tools, a.tools)
	return tools
}

// GetMode returns the agent's operation mode
func (a *baseAgent) GetMode() AgentMode {
	return a.mode
}

// GetMemory returns the agent's memory
func (a *baseAgent) GetMemory() Memory {
	return a.memory
}

// Run is implemented by specific agent types
func (a *baseAgent) Run(ctx context.Context, input string) (string, error) {
	return "", fmt.Errorf("Run() must be implemented by specific agent types")
}

// getToolDescriptions returns a formatted string of all tool descriptions
func (a *baseAgent) getToolDescriptions() string {
	var descriptions string
	for _, tool := range a.GetTools() {
		descriptions += fmt.Sprintf("Tool: %s\nDescription: %s\n\n", tool.Name(), tool.Description())
	}
	return descriptions
}
