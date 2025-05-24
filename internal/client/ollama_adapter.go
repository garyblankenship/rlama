package client

import (
	"context"

	"github.com/dontizi/rlama/internal/domain/agent"
)

// OllamaAdapter adapts OllamaClient to the agent's LLMClient interface
type OllamaAdapter struct {
	client *OllamaClient
	model  string
}

// NewOllamaAdapter creates a new adapter for OllamaClient
func NewOllamaAdapter(client *OllamaClient, model string) agent.LLMClient {
	return &OllamaAdapter{
		client: client,
		model:  model,
	}
}

// GenerateCompletion implements the agent.LLMClient interface
func (a *OllamaAdapter) GenerateCompletion(ctx context.Context, prompt string) (string, error) {
	// Call the existing OllamaClient
	return a.client.GenerateCompletion(a.model, prompt)
}

// GenerateEmbedding implements the agent.LLMClient interface
func (a *OllamaAdapter) GenerateEmbedding(ctx context.Context, text string) ([]float32, error) {
	// Call the existing OllamaClient with the same model used for completions
	return a.client.GenerateEmbedding(a.model, text)
}

// GenerateStructuredCompletion implements the agent.LLMClient interface
func (a *OllamaAdapter) GenerateStructuredCompletion(ctx context.Context, prompt string, schema map[string]interface{}) (string, error) {
	// Call the new structured completion method
	return a.client.GenerateStructuredCompletion(a.model, prompt, schema)
}
