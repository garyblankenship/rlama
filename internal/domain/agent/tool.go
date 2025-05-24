package agent

import (
	"context"
)

// Tool represents a tool that can be used by an agent
type Tool interface {
	Name() string
	Description() string
	Execute(ctx context.Context, input string) (string, error)
	// New methods for JSON schema support
	Schema() map[string]interface{}
	ExecuteWithParams(ctx context.Context, params map[string]interface{}) (string, error)
}

// ToolCall represents a structured tool call
type ToolCall struct {
	Name       string                 `json:"name"`
	Parameters map[string]interface{} `json:"parameters"`
}

// ToolResponse represents the agent's response with potential tool calls
type ToolResponse struct {
	Action     string                 `json:"action,omitempty"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	Response   string                 `json:"response,omitempty"`
	Reasoning  string                 `json:"reasoning,omitempty"`
	NeedsTool  bool                   `json:"needs_tool"`
}

// ToolExecution represents the result of a tool execution
type ToolExecution struct {
	Tool   string `json:"tool"`
	Result string `json:"result"`
	Error  string `json:"error,omitempty"`
}
